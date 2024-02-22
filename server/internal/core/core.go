package core

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	userentity "dalkak/internal/domain/user/object/entity"
	"dalkak/internal/infrastructure/database/dao"
	"dalkak/internal/infrastructure/eventbus"

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

	CreateUser(user *userentity.UserEntity) error
	FindUserByWalletAddress(walletAddress string) (*userentity.UserEntity, error)

	CreateUserMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error
	FindMediaTemp(userId, mediaType, prefix string) (*dao.MediaTempDao, error)
	// UpdateUserUploadMedia(userId string, findDto *dtos.MediaMeta, updateDto *dtos.UpdateUserUploadMediaDto) error
	// DeleteUserUploadMedia(userId string, dto *dtos.MediaMeta) error
}

type StorageManager interface {
	// GetHeadObject(key string) (*appdto.MediaHeadDto, error)
	// DeleteObject(key string) error
	CreatePresignedURL(mediaKey string, contentType string) (string, error)
}

type KeyManager interface {
	CreateSianature(sign string) (string, error)
	ParseTokenWithPublicKey(token string) (string, error)
}

type EventManager interface {
	Subscribe(eventType string, handler eventbus.EventHandler)
	Publish(event eventbus.Event)
}
