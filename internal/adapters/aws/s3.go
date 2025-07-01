package aws

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service interface {
	UploadFile(file *multipart.FileHeader, bucket string, key string) (string, error)
}

type s3Service struct {
	S3 *s3.S3
}

func NewS3Service(adapter *AwsAdapter) S3Service {
	return &s3Service{S3: s3.New(adapter.Session)}
}

func (s *s3Service) UploadFile(file *multipart.FileHeader, bucket string, key string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return "", err
	}

	_, err = s.S3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}
