package utils

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadImagePNGToS3(path string, fileBuffer []byte) error {
	accessToken := os.Getenv("AWS_ACCESS")
	secretToken := os.Getenv("AWS_SECRET")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	url := os.Getenv("AWS_URL")

	s3config := &aws.Config{
		Region:      aws.String("ru-1"),
		Credentials: credentials.NewStaticCredentials(accessToken, secretToken, ""),
		Endpoint:    aws.String(url),
	}

	s3session, err := session.NewSession(s3config)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(s3session)

	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path),
		Body:        bytes.NewReader(fileBuffer),
		ContentType: aws.String("image/png"),
	}

	_, err = uploader.UploadWithContext(context.Background(), uploadInput)

	return err
}
