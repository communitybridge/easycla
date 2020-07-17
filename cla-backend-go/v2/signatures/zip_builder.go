package signatures

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/juju/zip"

	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// constants
const (
	LocalFolder        = "/mnt/storage"
	ICLA               = "icla"
	CCLA               = "ccla"
	ParallelDownloader = 100
)

type Zipper struct {
	s3         *s3.S3
	bucketName string
}

type ZipBuilder interface {
	BuildICLAZip(claGroupID string) error
	BuildCCLAZip(claGroupID string) error
}

func NewZipBuilder(awsSession *session.Session, bucketName string) ZipBuilder {
	return &Zipper{
		s3:         s3.New(awsSession),
		bucketName: bucketName,
	}
}

func s3ZipFilepath(claType string, claGroupID string) string {
	return fmt.Sprintf("contract-group/%s/%s.zip", claGroupID, claType)
}

func localZipFilepath(claType string, claGroupID string) string {
	return fmt.Sprintf("%s/%s_%s.zip", LocalFolder, claGroupID, claType)
}

func s3ZipPrefix(claType string, claGroupID string) string {
	return fmt.Sprintf("contract-group/%s/%s/", claGroupID, claType)
}

func (z *Zipper) BuildICLAZip(claGroupID string) error {
	return z.buildZip(ICLA, claGroupID)
}

func (z *Zipper) BuildCCLAZip(claGroupID string) error {
	return z.buildZip(CCLA, claGroupID)
}

func (z *Zipper) buildZip(claType string, claGroupID string) error {
	f := logrus.Fields{"cla_group_id": claGroupID, "cla_type": claType}
	log.WithFields(f).Debug("syncing zip file for cla group")
	err := z.syncFile(claType, claGroupID)
	if err != nil {
		log.WithFields(f).Errorf("syncing zip file for cla group failed. error = %v", err)
		return err
	}

	log.WithFields(f).Debug("syncing zip file for cla group")
	localZipFile := localZipFilepath(claType, claGroupID)
	_, err = os.Stat(localZipFile)
	var newFile bool
	if err != nil {
		if os.IsNotExist(err) {
			newFile = true
			log.WithFields(f).Debug("zip file not exist")
		} else {
			return err
		}
	}
	var files *utils.StringSet
	if !newFile {
		log.WithFields(f).Debug("getting list of all zip files in list")
		files, err = getZipFiles(localZipFile)
		if err != nil {
			return err
		}
	}
	log.WithFields(f).Debug("getting zip writer")
	writer, zfile, err := getZipWriter(newFile, claType, claGroupID)
	if err != nil {
		return err
	}
	var zipUpdated bool
	log.WithFields(f).Debug("getting s3 files")
	downloaderInputChan := make(chan *DownloadFileInput)
	downloaderOutputChan := make(chan *FileContent)
	var wg sync.WaitGroup
	wg.Add(ParallelDownloader)
	for i := 1; i <= ParallelDownloader; i++ {
		go z.downloader(&wg, downloaderInputChan, downloaderOutputChan)
	}
	go func() {
		wg.Wait()
		close(downloaderOutputChan)
	}()
	go func() {
		err = z.s3.ListObjectsPages(&s3.ListObjectsInput{
			Bucket: aws.String(z.bucketName),
			Prefix: aws.String(s3ZipPrefix(claType, claGroupID)),
		}, func(output *s3.ListObjectsOutput, b bool) bool {
			for _, obj := range output.Contents {
				key := utils.StringValue(obj.Key)
				tmp := strings.Split(key, "/")
				if len(tmp) != 5 {
					continue
				}
				filename := tmp[4]
				if !newFile && files.Include(filename) {
					// skip files which are already present in zip
					continue
				}
				downloaderInputChan <- &DownloadFileInput{
					filename: filename,
					key:      obj.Key,
				}
			}
			return true
		})
		close(downloaderInputChan)
	}()
	zipUpdated = writeFileToZip(writer, downloaderOutputChan)
	if err != nil {
		zfile.Close()
		return err
	}
	if zipUpdated {
		// we should only close the writer when we have have written some file
		writer.Close()
	}
	zfile.Close()
	if zipUpdated {
		remoteZipFileKey := s3ZipFilepath(claType, claGroupID)
		log.Debugf("Uploading zip file %s", remoteZipFileKey)
		err := z.UploadFile(localZipFile, remoteZipFileKey)
		if err != nil {
			log.Warnf("Uploading zip file %s failed. error = %s", remoteZipFileKey, err.Error())
			return err
		}
		log.Debugf("Uploaded zip file %s", remoteZipFileKey)
	}
	return nil
}

type FileContent struct {
	buff     *aws.WriteAtBuffer
	filename string
}

type DownloadFileInput struct {
	filename string
	key      *string
}

func writeFileToZip(writer *zip.Writer, filesInput chan *FileContent) bool {
	var zipUpdated bool
	for fileContent := range filesInput {
		filename := fileContent.filename
		buff := fileContent.buff
		log.Debugf("Adding file : %s to zip", filename)
		header := &zip.FileHeader{
			Name:   filename,
			Method: zip.Deflate,
		}
		header.SetMode(0644)
		f, err := writer.CreateHeader(header)
		if err != nil {
			log.WithField("file", filename).Error("unable to write file header in zip")
			continue
		}
		_, err = f.Write(buff.Bytes())
		if err != nil {
			log.WithField("file", filename).Error("unable to write file data in zip")
			continue
		}
		zipUpdated = true
	}
	return zipUpdated
}

func (z *Zipper) downloader(wg *sync.WaitGroup, inputChan chan *DownloadFileInput, outputChan chan *FileContent) {
	defer wg.Done()
	for in := range inputChan {
		log.Debugf("Downloading file : %s", in.filename)
		buff := &aws.WriteAtBuffer{}
		downloader := s3manager.NewDownloaderWithClient(z.s3)
		_, err := downloader.Download(buff,
			&s3.GetObjectInput{
				Bucket: aws.String(z.bucketName),
				Key:    in.key,
			})
		if err != nil {
			log.WithField("key", utils.StringValue(in.key)).Error("unable to download file from s3", err)
			continue
		}
		outputChan <- &FileContent{
			buff:     buff,
			filename: in.filename,
		}
	}
}

func getZipFiles(zipFile string) (*utils.StringSet, error) {
	files := utils.NewStringSet()
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, file := range r.File {
		files.Add(file.Name)
	}
	return files, nil
}

func getZipWriter(newFile bool, claType string, claGroupID string) (*zip.Writer, *os.File, error) {
	zipFile := localZipFilepath(claType, claGroupID)
	var zfile *os.File
	var err error
	var writer *zip.Writer
	if newFile {
		zfile, err = os.OpenFile(zipFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, nil, err
		}
		writer = zip.NewWriter(zfile)
	} else {
		// open existing zip file
		zfile, err = os.OpenFile(zipFile, os.O_RDWR, 0644)
		if err != nil {
			return nil, nil, err
		}
		// seek to end of file
		size, err := zfile.Seek(0, 2)
		if err != nil {
			zfile.Close()
			return nil, nil, err
		}
		r, err := zip.NewReader(zfile, size)
		if err != nil {
			zfile.Close()
			return nil, nil, err
		}
		writer = r.Append(zfile)
	}
	return writer, zfile, nil
}

func (z *Zipper) syncFile(claType string, claGroupID string) error {
	var localFileExist bool
	var remoteFileExist bool
	localFilePath := localZipFilepath(claType, claGroupID)
	remoteFilePath := s3ZipFilepath(claType, claGroupID)
	var syncNeeded bool
	log.Debugf("checking if localfile exist or not. file = %s", localFilePath)
	info, err := os.Stat(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			localFileExist = false
		} else {
			return err
		}
	} else {
		localFileExist = true
	}
	log.Debugf("local file exist: %v", localFileExist)
	log.Debugf("Checking if remote file exist. file =  %v", remoteFilePath)
	// check if no s3 file exist or not
	resp, err := z.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(z.bucketName),
		Key:    aws.String(remoteFilePath),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			remoteFileExist = false
		}
	} else {
		remoteFileExist = true
	}
	log.Debugf("remote file exist: %v", remoteFileExist)
	if remoteFileExist {
		if localFileExist {
			if resp.ContentLength != nil && *resp.ContentLength != info.Size() {
				syncNeeded = true
				log.WithField("remote_file_size", *resp.ContentLength).
					WithField("local_file_size", info.Size()).
					Debug("content of remote and local file is not same")
			} else {
				log.Debug("content of remote and local file is same")
			}
		} else {
			syncNeeded = true
		}
	} else {
		if localFileExist {
			os.Remove(localFilePath)
			localFileExist = false
		}
	}
	if syncNeeded {
		err = z.downloadZipFile(claType, claGroupID)
		if err != nil {
			log.Error("Downloading zip file failed", err)
			return err
		}
	}
	return nil
}

func (z *Zipper) downloadZipFile(claType string, claGroupID string) error {
	localFile := localZipFilepath(claType, claGroupID)
	remoteFileKey := s3ZipFilepath(claType, claGroupID)
	log.Debugf("Downloading zip file %s", remoteFileKey)
	file, err := os.Create(localFile)
	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloaderWithClient(z.s3)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(z.bucketName),
			Key:    aws.String(remoteFileKey),
		})
	if err != nil {
		return err
	}
	log.Debugf("Downloading zip file %s completed", remoteFileKey)
	return nil
}

func (z *Zipper) UploadFile(localFilename, s3ZipFile string) error {
	uploader := s3manager.NewUploaderWithClient(z.s3)

	// Open file
	f, err := os.Open(localFilename)
	if err != nil {
		log.Debugf("failed to open file %q, %v", localFilename, err)
		return err
	}
	defer f.Close()

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(z.bucketName),
		Key:    aws.String(s3ZipFile),
		Body:   f,
	})

	//in case it fails to upload
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return err
	}
	return nil
}
