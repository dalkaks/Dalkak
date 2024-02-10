package app

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/generateutils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client     *s3.Client
	bucket     string
	staticLink string
}

var _ interfaces.Storage = (*Storage)(nil)

func NewStorage(ctx context.Context, mode string, staticLink string) (*Storage, error) {
	var cfg aws.Config
	var err error

	if mode == "LOCAL" {
		cfg, err = awsConfig.LoadDefaultConfig(ctx, awsConfig.WithSharedConfigProfile("dalkak"))
	} else {
		cfg, err = awsConfig.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	storageClient := s3.NewFromConfig(cfg)
	var bucket string
	if mode == "PROD" {
		bucket = "dalkak-prod"
	} else {
		bucket = "dalkak-dev"
	}

	return &Storage{client: storageClient, bucket: bucket, staticLink: staticLink}, nil
}

func (storage *Storage) CreatePresignedURL(mediaType dtos.MediaType, ext string) (*dtos.MediaMeta, error) {
	expires := 30 * time.Minute
	presigner := s3.NewPresignClient(storage.client, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	contentType := fmt.Sprintf("%s/%s", mediaType, ext)
	id, err := storage.generateMediaId(mediaType)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("temp/%s/%s.%s", mediaType, id, ext)

	presignedURL, err := presigner.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(storage.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create presigned url",
		}
	}

	return &dtos.MediaMeta{
		ID:          id,
		Extension:   ext,
		ContentType: contentType,
		URL:         presignedURL.URL,
	}, nil
}

func (storage *Storage) generateMediaId(mediaType dtos.MediaType) (string, error) {
	uuid := generateutils.GenerateUUID()
	key := fmt.Sprintf("temp/%s/%s", mediaType.String(), uuid)

	_, err := storage.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})

	var noSuchKeyErr *types.NoSuchKey
	var notFoundErr *types.NotFound
	if err != nil && (errors.As(err, &noSuchKeyErr) || errors.As(err, &notFoundErr)) {
		return uuid, nil
	}
	// Todo: error handling(중복 또는 에러)
	return "", &dtos.AppError{
		Code:    http.StatusInternalServerError,
		Message: "Failed to generate media id",
	}
}
