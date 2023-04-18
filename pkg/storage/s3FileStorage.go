// storage/s3_file_storage.go
package storage

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3FileStorage struct {
	s3API      s3iface.S3API
	bucket     string
	downloader *s3manager.Downloader
}

func NewS3FileStorage(region, bucket string) (*S3FileStorage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	s3API := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	return &S3FileStorage{
		s3API:      s3API,
		bucket:     bucket,
		downloader: downloader,
	}, nil
}

func (s *S3FileStorage) GetFile(path string) (io.ReadSeekCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	result, err := s.s3API.GetObject(input)
	if err != nil {
		return nil, err
	}

	return toReadSeekCloser(result.Body), nil
}
