package signatures

import (
	"bytes"
	"fmt"
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
	// get zip file from s3
	buff, err := z.getZipFileFromS3(claType, claGroupID)
	if err != nil {
		return err
	}
	var files *utils.StringSet
	if len(buff.Bytes()) != 0 {
		// read files already present in zip
		log.Debug("reading files present in zip")
		files, err = getZipFiles(buff)
		if err != nil {
			return err
		}
	}
	log.WithFields(f).Debug("getting zip writer")
	writer, err := getZipWriter(buff)
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
				if files != nil && files.Include(filename) {
					// skip files which are already present in zip
					log.Debugf("file %s already present in zip", filename)
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
		return err
	}
	writer.Close()
	if zipUpdated {
		remoteZipFileKey := s3ZipFilepath(claType, claGroupID)
		log.Debugf("Uploading zip file %s", remoteZipFileKey)
		err := z.UploadFile(buff, remoteZipFileKey)
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

func getZipFiles(buff *bytes.Buffer) (*utils.StringSet, error) {
	reader := bytes.NewReader(buff.Bytes())
	files := utils.NewStringSet()
	r, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return nil, err
	}
	for _, file := range r.File {
		files.Add(file.Name)
	}
	return files, nil
}

func getZipWriter(buff *bytes.Buffer) (*zip.Writer, error) {
	var writer *zip.Writer
	if len(buff.Bytes()) == 0 {
		writer = zip.NewWriter(buff)
	} else {
		reader := bytes.NewReader(buff.Bytes())
		r, err := zip.NewReader(reader, reader.Size())
		if err != nil {
			return nil, err
		}
		writer = r.Append(buff)
	}
	return writer, nil
}

func (z *Zipper) getZipFileFromS3(claType string, claGroupID string) (*bytes.Buffer, error) {
	var buff aws.WriteAtBuffer
	remoteFileKey := s3ZipFilepath(claType, claGroupID)
	_, err := z.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(z.bucketName),
		Key:    aws.String(remoteFileKey),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			log.Debugf("zip file %s does not exist on s3", remoteFileKey)
			return bytes.NewBuffer(buff.Bytes()), nil
		}
		return nil, err
	}
	log.Debugf("Downloading zip file %s", remoteFileKey)

	downloader := s3manager.NewDownloaderWithClient(z.s3)

	_, err = downloader.Download(&buff,
		&s3.GetObjectInput{
			Bucket: aws.String(z.bucketName),
			Key:    aws.String(remoteFileKey),
		})
	if err != nil {
		return nil, err
	}
	log.Debugf("Downloading zip file %s completed", remoteFileKey)
	return bytes.NewBuffer(buff.Bytes()), nil
}

func (z *Zipper) UploadFile(localFileContent *bytes.Buffer, s3ZipFile string) error {
	uploader := s3manager.NewUploaderWithClient(z.s3)
	// Upload the file to S3.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(z.bucketName),
		Key:    aws.String(s3ZipFile),
		Body:   localFileContent,
	})

	//in case it fails to upload
	if err != nil {
		log.Warnf("failed to upload file %s. error = %v", s3ZipFile, err)
		return err
	}
	return nil
}
