package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"ps-beli-mang/configs"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/helper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type ImageService interface {
	UploadImage(multipart.File, *multipart.FileHeader) (string, error)
}

type imageService struct {
	s3Client   *s3.Client
	bucketName string
	cfg        *configs.MainConfig
}

func NewImageService(cfg *configs.MainConfig) ImageService {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3.ID, cfg.S3.SecretKey, "")),
	)
	if err != nil {
		helper.PanicIfError(err, "Error when reading env config.")
	}

	s3Client := s3.NewFromConfig(awsConfig)

	return &imageService{
		s3Client:   s3Client,
		bucketName: cfg.S3.BucketName,
		cfg:        cfg,
	}
}

func (s *imageService) UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	if fileHeader.Size > 2*1024*1024 || fileHeader.Size < 10*1024 {
		return "", errs.NewErrBadRequest("file size must be between 10KB and 2MB")
	}

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/jpg" {
		return "", errs.NewErrBadRequest("invalid file type")
	}

	fileName := uuid.New().String() + ".jpeg"

	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	})
	helper.PanicIfError(err, "Failed to upload to S3.")

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.cfg.S3.Region, fileName)
	return fileUrl, nil
}
