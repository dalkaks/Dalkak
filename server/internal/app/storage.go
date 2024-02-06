package app

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/generateutils"
	"errors"
	"fmt"
	"strings"

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
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
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

func (storage *Storage) Upload(media *dtos.MediaDto, path string) (*dtos.MediaMeta, error) {
	contentTypeParts := strings.Split(media.Meta.ContentType, "/")
	contentTypePath := ""
	if len(contentTypeParts) > 0 {
		contentTypePath = contentTypeParts[0]
	}

	uuid := generateutils.GenerateUUID()
	key := fmt.Sprintf("%s/%s/%s.%s", path, contentTypePath, uuid, media.Meta.Extension)

	_, err := storage.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var noSuchKeyErr *types.NoSuchKey
		var notFoundErr *types.NotFound
		if errors.As(err, &noSuchKeyErr) || errors.As(err, &notFoundErr) {
			_, err := storage.client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket:      aws.String(storage.bucket),
				Key:         aws.String(key),
				Body:        media.File,
				ContentType: aws.String(media.Meta.ContentType),
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("file already exists")
	}

	return &dtos.MediaMeta{
		ID:          uuid,
		Extension:   media.Meta.Extension,
		ContentType: media.Meta.ContentType,
		URL: storage.staticLink + key,
	}, nil
}
