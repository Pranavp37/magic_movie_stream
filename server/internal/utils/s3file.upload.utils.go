package utils

import (
	// "github.com/Pranavp37/magic_movie_stream/internal/config"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/Pranavp37/magic_movie_stream/internal/configs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Uploader struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
	region     string
}

func NewS3Uploader(ctx context.Context) (*S3Uploader, error) {

	bucketName := configs.LoadConfig().BucketName

	region := configs.LoadConfig().S3_Region
	awsAccessKeyID := configs.LoadConfig().AWS_ACCESS_KEY_ID
	awsSecretAccessKey := configs.LoadConfig().AWS_SECRET_ACCESS_KEY

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
		awsAccessKeyID,
		awsSecretAccessKey,
		"",
	)))

	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(s3Client)

	return &S3Uploader{
		client:     s3Client,
		uploader:   uploader,
		bucketName: bucketName,
		region:     region,
	}, nil
}

func (s *S3Uploader) FileUploader(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	key := "profiles/" + uuid.New().String() + ext

	_, err = s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   io.Reader(file),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)

	return url, nil
}
