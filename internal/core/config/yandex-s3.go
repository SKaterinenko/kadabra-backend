package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3Client(cfg *Config) (*S3Client, error) {
	cfgS3, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ru-central1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(cfgS3, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
		o.UsePathStyle = true
	})

	return &S3Client{
		client:     client,
		bucketName: cfg.S3Bucket,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Генерируем уникальное имя файла
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Загружаем файл
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		//ACL: types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Формируем URL
	url := fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", s.bucketName, filename)
	return url, nil
}

// Удаление файла
func (s *S3Client) DeleteFile(ctx context.Context, filename string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	})
	return err
}
