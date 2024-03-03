package storage

import (
	"context"
	storagedto "dalkak/internal/infrastructure/storage/type"
	parseutil "dalkak/pkg/utils/parse"
	responseutil "dalkak/pkg/utils/response"
	"errors"
	"log"
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

const presignExpireMinutes = 10 * time.Minute

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

func (storage *Storage) CreatePresignedURL(mediaKey string, contentType string) (string, error) {
	presigner := s3.NewPresignClient(storage.client, func(o *s3.PresignOptions) {
		o.Expires = presignExpireMinutes
	})

	presignedURL, err := presigner.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(storage.bucket),
		Key:         aws.String(mediaKey),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}

	return presignedURL.URL, nil
}

func (storage *Storage) CopyObject(srcUrl, descUrl string) error {
	srcKey := parseutil.ConvertStaticLinkToKey(storage.staticLink, srcUrl)
	descKey := parseutil.ConvertStaticLinkToKey(storage.staticLink, descUrl)
	log.Println(srcKey, descKey)

	_, err := storage.client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(storage.bucket),
		CopySource: aws.String(storage.bucket + "/" + srcKey),
		Key:        aws.String(descKey),
	})
	// todo log
	if err != nil {
		log.Println(err)
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}
	return nil
}

func (storage *Storage) GetHeadObject(key string) (*storagedto.MediaHeadDto, error) {
	headObjectOutput, err := storage.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if errors.Is(err, &types.NoSuchKey{}) || errors.Is(err, &types.NotFound{}) {
			return nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgDataNotFound, err)
		}
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}
	return &storagedto.MediaHeadDto{
		ContentType: *headObjectOutput.ContentType,
		Length:      *headObjectOutput.ContentLength,
	}, nil
}

func (storage *Storage) DeleteObject(key string) error {
	_, err := storage.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}
	return nil
}
