package adapters

import (
	"export-service/internal/core/ports"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"os"
)

type S3Uploader struct {
	client *s3.S3
	logger *zap.Logger
	bucket string
	folder string
}

var _ ports.Uploader = (*S3Uploader)(nil)

func NewS3Uploader(accessKey, secretKey, endpoint, region, bucket, folder string, logger *zap.Logger) *S3Uploader {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		logger.Fatal("Failed to create S3 session", zap.Error(err))
	}

	return &S3Uploader{
		client: s3.New(sess),
		logger: logger,
		bucket: bucket,
		folder: folder,
	}
}

func (s *S3Uploader) Upload(fileName, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		s.logger.Error("Failed to open file", zap.Error(err), zap.String("path", path))
		return "", err
	}
	defer file.Close()

	key := fmt.Sprintf("%s/%s", s.folder, fileName)
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		s.logger.Error("Failed to upload file", zap.Error(err), zap.String("path", path))
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.bhs.io.cloud.ovh.net/%s", s.bucket, key)

	s.logger.Info("Successfully uploaded file", zap.String("url", url), zap.String("path", path))

	return url, nil
}
