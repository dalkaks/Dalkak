package core

import (
	userobject "dalkak/internal/domain/user/object"
	"dalkak/internal/infrastructure/eventbus"
	appdto "dalkak/pkg/dto/app"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Infra struct {
	Database     DatabaseManager
	Storage      StorageManager
	Keymanager   KeyManager
	EventManager EventManager
}

type DatabaseManager interface {
	GetClient() *dynamodb.Client
	GetTable() string

	CreateUser(user *userobject.UserEntity) error
	FindUserByWalletAddress(walletAddress string) (*userobject.UserEntity, error)

	// CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error
	// FindUserUploadMedia(userId string, dto *dtos.FindUserUploadMediaDto) (*dtos.MediaMeta, error)
	// UpdateUserUploadMedia(userId string, findDto *dtos.MediaMeta, updateDto *dtos.UpdateUserUploadMediaDto) error
	// DeleteUserUploadMedia(userId string, dto *dtos.MediaMeta) error
}

type StorageManager interface {
	GetHeadObject(key string) (*appdto.MediaHeadDto, error)
	DeleteObject(key string) error
	CreatePresignedURL(userId string, dto *appdto.UploadMediaDto) (*appdto.MediaMeta, string, error)
	ConvertStaticLinkToKey(url string) (string, error)
}

type KeyManager interface {
	CreateSianature(sign string) (string, error)
	ParseTokenWithPublicKey(token string) (string, error)
}

type EventManager interface {
	Subscribe(eventType string, handler eventbus.EventHandler)
	Publish(event eventbus.Event)
}