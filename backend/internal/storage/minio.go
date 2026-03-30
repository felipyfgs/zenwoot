package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/felipyfgs/zenwoot/backend/internal/config"
	"github.com/felipyfgs/zenwoot/backend/internal/logger"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(cfg config.MinIOConfig) (*S3Storage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
		awsconfig.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, fmt.Errorf("aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s", cfg.Endpoint))
		o.UsePathStyle = true
	})

	return &S3Storage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}

	objectName := fmt.Sprintf("%s/%d/%s", folder, time.Now().Unix(), file.Filename)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("s3 put: %w", err)
	}

	logger.Info().Str("object", objectName).Msg("file uploaded")
	return objectName, nil
}

func (s *S3Storage) GetURL(ctx context.Context, objectName string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	presigned, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
	}, s3.WithPresignExpires(24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}
	return presigned.URL, nil
}

func (s *S3Storage) GetObject(ctx context.Context, objectName string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}
	return resp.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, objectName string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return fmt.Errorf("s3 delete: %w", err)
	}
	return nil
}
