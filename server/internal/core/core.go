package core

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	userentity "dalkak/internal/domain/user/object/entity"
	"dalkak/internal/infrastructure/database/dao"
	"dalkak/internal/infrastructure/eventbus"
	keytype "dalkak/internal/infrastructure/key/type"
	storagedto "dalkak/internal/infrastructure/storage/type"

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

	GetTransactionID() (*dao.TransactionDao, error)

	CreateUser(user *userentity.UserEntity) error
	FindUserByWalletAddress(walletAddress string) (*dao.UserDao, error)

	CreateMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error
	FindMediaTemp(userId, mediaType, prefix string) (*dao.MediaTempDao, error)
	UpdateMediaTempConfirm(userId string, mediaTempUpdate *mediaaggregate.MediaTempUpdate) error
	DeleteMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error

	CreateBoard(txId string, board *boardaggregate.BoardAggregate, order *orderaggregate.OrderAggregate, imageResource, videoResource *mediavalueobject.MediaResource) error
	FindBoardByUserId(userId string, status *boardentity.BoardStatus, pageDao *dao.RequestPageDao) ([]*dao.BoardDao, *dao.ResponsePageDao, error)
}

type StorageManager interface {
	CreatePresignedURL(mediaKey string, contentType string) (string, error)
	GetHeadObject(key string) (*storagedto.MediaHeadDto, error)
	DeleteObject(key string) error
}

type KeyManager interface {
	CreateSianature(sign string) (string, error)
	ParseTokenWithPublicKey(token string, tokenType keytype.TokenType) (string, error)
	GetDomain() string
	GetMode() string
}

type EventManager interface {
	Subscribe(eventType string, handler eventbus.EventHandler)
	Publish(event eventbus.Event)
}
