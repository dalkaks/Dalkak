package storage

import (
	"context"
	appdto "dalkak/pkg/dto/app"
	responseutil "dalkak/pkg/utils/response"
	generateutil "dalkak/pkg/utils/generate"
	"errors"
	"strings"
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

const prefixExpireMinutes = 10

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

func (storage *Storage) GetHeadObject(key string) (*appdto.MediaHeadDto, error) {
	headObjectOutput, err := storage.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if errors.Is(err, &types.NoSuchKey{}) || errors.Is(err, &types.NotFound{}) {
			return nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgStorageNoSuchKey, err)
		}
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}
	return &appdto.MediaHeadDto{
		Key:         key,
		ContentType: *headObjectOutput.ContentType,
		Length:      *headObjectOutput.ContentLength,
		URL:         storage.convertKeyToStaticLink(key),
		MetaUserId:  headObjectOutput.Metadata["userid"],
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

func (storage *Storage) CreatePresignedURL(userId string, dto *appdto.UploadMediaDto) (*appdto.MediaMeta, string, error) {
	mediaType := dto.MediaType.String()
	expires := prefixExpireMinutes * time.Minute
	presigner := s3.NewPresignClient(storage.client, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	contentType := GenerateContentType(mediaType, dto.Ext)
	id, err := storage.generateMediaId(dto)
	if err != nil {
		return nil, "", err
	}

	key := GenerateTempKey(dto.Prefix, mediaType, id, dto.Ext)

	presignedURL, err := presigner.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(storage.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"userid": userId,
		},
	})
	if err != nil {
		return nil, "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
	}
	storageUrl := storage.convertKeyToStaticLink(key)

	return &appdto.MediaMeta{
		ID:          id,
		Prefix:      dto.Prefix,
		Extension:   dto.Ext,
		ContentType: contentType,
		URL:         storageUrl,
	}, presignedURL.URL, nil
}

func (storage *Storage) ConvertStaticLinkToKey(url string) (string, error) {
	if url == "" || !strings.HasPrefix(url, storage.staticLink) {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgStorageInvalidURL, errors.New("invalid url"))
	}
	return url[len(storage.staticLink):], nil
}

func (storage *Storage) generateMediaId(dto *appdto.UploadMediaDto) (string, error) {
	uuid := generateutil.GenerateUUID()
	key := GenerateMediaPath(dto.Prefix, dto.MediaType.String(), uuid)

	_, err := storage.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(storage.bucket),
		Key:    aws.String(key),
	})

	var noSuchKeyErr *types.NoSuchKey
	var notFoundErr *types.NotFound
	if err != nil && (errors.As(err, &noSuchKeyErr) || errors.As(err, &notFoundErr)) {
		return uuid, nil
	}
	// Todo: temp 파일 체크, error handling(중복 또는 에러)
	return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgStorageInternal, err)
}

func (storage *Storage) convertKeyToStaticLink(key string) string {
	return storage.staticLink + key
}


func GenerateTempKey(prefix string, mediaType string, id string, ext string) string {
	return "temp/" + prefix + "/" + mediaType + "/" + id + "." + ext
}

func GenerateMediaPath(prefix string, mediaType string, id string) string {
	return prefix + "/" + mediaType + "/" + id
}

func GenerateContentType(mediaType string, ext string) string {
	return mediaType + "/" + ext
}
