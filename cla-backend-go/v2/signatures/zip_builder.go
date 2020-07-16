package signatures

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/utils"

	"github.com/juju/zip"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// constants
const (
	LocalFolder = "/mnt/storage"
	ICLA        = "icla"
	CCLA        = "ccla"
)

type Zipper struct {
	s3         *s3.S3
	bucketName string
}

type ZipBuilder interface {
	BuildICLAZip(claGroupID string) error
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
	f := logrus.Fields{"cla_group_id": claGroupID, "cla_type": ICLA}
	log.WithFields(f).Debug("syncing zip file for cla group")
	err := z.syncFile(ICLA, claGroupID)
	if err != nil {
		log.WithFields(f).Errorf("syncing zip file for cla group failed. error = %v", err)
		return err
	}

	log.WithFields(f).Debug("syncing zip file for cla group")
	zipFile := localZipFilepath(ICLA, claGroupID)
	_, err = os.Stat(zipFile)
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
		files, err = getZipFiles(zipFile)
		if err != nil {
			return err
		}
	}
	log.WithFields(f).Debug("getting zip writer")
	writer, err := getZipWriter(newFile, ICLA, claGroupID)
	if err != nil {
		return err
	}
	var zipUpdated bool
	log.WithFields(f).Debug("getting s3 files")
	err = z.s3.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(z.bucketName),
		Prefix: aws.String(s3ZipPrefix(ICLA, claGroupID)),
	}, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, obj := range output.Contents {
			key := utils.StringValue(obj.Key)
			log.Debugf("filename : %s", key)
			tmp := strings.Split(key, "/")
			if len(tmp) != 5 {
				continue
			}
			filename := tmp[4]
			if !newFile && files.Include(filename) {
				// skip files which are already present in zip
				continue
			}
			log.Debugf("Downloading file : %s", filename)
			buff := &aws.WriteAtBuffer{}
			downloader := s3manager.NewDownloaderWithClient(z.s3)
			_, err = downloader.Download(buff,
				&s3.GetObjectInput{
					Bucket: aws.String(z.bucketName),
					Key:    obj.Key,
				})
			if err != nil {
				log.WithField("file", key).Error("unable to add file to zip")
				continue
			}
			header := &zip.FileHeader{
				Name:   filename,
				Method: zip.Deflate,
			}
			header.SetMode(0644)
			f, err := writer.CreateHeader(header)
			if err != nil {
				log.WithField("file", key).Error("unable to write file header in zip")
				continue
			}
			_, err = f.Write(buff.Bytes())
			if err != nil {
				log.WithField("file", key).Error("unable to write file data in zip")
				continue
			}
			zipUpdated = true
		}
		return true
	})
	if err != nil {
		return err
	}
	if zipUpdated {

	}
	return nil
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

func getZipWriter(newFile bool, claType string, claGroupID string) (*zip.Writer, error) {
	zipFile := localZipFilepath(claType, claGroupID)
	var writer *zip.Writer
	if newFile {
		zfile, err := os.OpenFile(zipFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		writer = zip.NewWriter(zfile)
	} else {
		zfile, err := os.OpenFile(zipFile, os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		size, err := zfile.Seek(0, 2)
		if err != nil {
			return nil, err
		}
		r, err := zip.NewReader(zfile, size)
		if err != nil {
			return nil, err
		}
		writer = r.Append(zfile)
	}
	return writer, nil
}

func (z *Zipper) syncFile(claGroupID string, claType string) error {
	var localFileExist bool
	var remoteFileExist bool
	var syncNeeded bool
	log.Debug("checking if localfile exist or not")
	info, err := os.Stat(localZipFilepath(claType, claGroupID))
	if err != nil {
		if os.IsNotExist(err) {
		} else {
			return err
		}
	} else {
		localFileExist = true
	}
	// check if no s3 file exist or not
	resp, err := z.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(z.bucketName),
		Key:    aws.String(s3ZipFilepath(claType, claGroupID)),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			remoteFileExist = false
		}
	} else {
		remoteFileExist = true
	}
	if remoteFileExist {
		if localFileExist {
			if resp.ContentLength != nil && *resp.ContentLength != info.Size() {
				syncNeeded = true
			}
		} else {
			syncNeeded = true
		}
	} else {
		localFileExist = false
		os.Remove(localZipFilepath(claType, claGroupID))
	}
	if syncNeeded {
		err = z.syncS3File(claType, claGroupID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *Zipper) syncS3File(claType string, claGroupID string) error {
	localFile := localZipFilepath(claType, claGroupID)
	file, err := os.Create(localFile)
	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloaderWithClient(z.s3)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(z.bucketName),
			Key:    aws.String(s3ZipFilepath(claType, claGroupID)),
		})
	if err != nil {
		return err
	}
	logging.WithField("localfile", localFile).Debug("synced s3 file")
	return nil
}

func (z *Zipper) UploadFile(localFilename, s3ZipFile string) {
	uploader := s3manager.NewUploaderWithClient(z.s3)

	//open the file
	f, err := os.Open(localFilename)
	if err != nil {
		log.Debugf("failed to open file %q, %v", filename, err)
		return
	}
	//defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(myKey),
		Body:   f,
	})

	//in case it fails to upload
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return
	}
}
