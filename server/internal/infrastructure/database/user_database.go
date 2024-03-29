package database

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	userentity "dalkak/internal/domain/user/object/entity"
	"dalkak/internal/infrastructure/database/dao"
	responseutil "dalkak/pkg/utils/response"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const UserDataType = "User"

func GenerateUserDataPk(walletAddress string) string {
	return UserDataType + `#` + walletAddress
}

func GenerateUserBoardImageDataSk(prefix string, mediaType string) string {
	return `Media#` + prefix + `#` + mediaType
}

type UserData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	WalletAddress string
}

type UserMediaData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id          string
	Prefix      string
	Extension   string
	ContentType string
	Url         string
	IsConfirm   bool
}

func (repo *Database) CreateUser(user *userentity.UserEntity) error {
	pk := GenerateUserDataPk(user.WalletAddress)
	newUser := &UserData{
		Pk:         pk,
		Sk:         pk,
		EntityType: UserDataType,
		Timestamp:  user.Timestamp,

		WalletAddress: user.WalletAddress,
	}

	err := repo.PutDynamoDBItem(newUser)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Database) FindUserByWalletAddress(walletAddress string) (*dao.UserDao, error) {
	pk := GenerateUserDataPk(walletAddress)
	var userToFind *UserData

	keyCond := expression.Key("Pk").Equal(expression.Value(pk)).
		And(expression.Key("Sk").Equal(expression.Value(pk)))
	expr, err := GenerateQueryExpression(keyCond, nil)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	err = repo.QuerySingleItem(expr, &userToFind)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}
	if userToFind == nil {
		return nil, nil
	}

	return &dao.UserDao{
		WalletAddress: userToFind.WalletAddress,
		Timestamp:     userToFind.Timestamp,
	}, nil
}

func (repo *Database) CreateMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error {
	prefix := mediaTemp.MediaResource.Prefix.String()
	mediaType := mediaTemp.MediaResource.ContentType.ConvertToMediaType()

	sk := GenerateUserBoardImageDataSk(prefix, mediaType)

	newUploadMedia := &UserMediaData{
		Pk:         GenerateUserDataPk(userId),
		Sk:         sk,
		EntityType: sk,
		Timestamp:  mediaTemp.MediaEntity.Timestamp,

		Id:          mediaTemp.MediaEntity.Id,
		Prefix:      prefix,
		Extension:   mediaTemp.MediaResource.GetExtension(),
		ContentType: mediaTemp.MediaResource.ContentType.String(),
		Url:         mediaTemp.MediaUrl.AccessUrl,
		IsConfirm:   mediaTemp.MediaEntity.IsConfirm,
	}

	err := repo.PutDynamoDBItem(newUploadMedia)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Database) FindMediaTemp(userId, mediaType, prefix string) (*dao.MediaTempDao, error) {
	sk := GenerateUserBoardImageDataSk(prefix, mediaType)
	var mediaToFind *UserMediaData

	keyCond := expression.Key("Pk").Equal(expression.Value(GenerateUserDataPk(userId))).
		And(expression.Key("Sk").Equal(expression.Value(sk)))
	expr, err := GenerateQueryExpression(keyCond, nil)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	err = repo.QuerySingleItem(expr, &mediaToFind)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}
	if mediaToFind == nil {
		return nil, nil
	}

	return &dao.MediaTempDao{
		Id:          mediaToFind.Id,
		Prefix:      mediaToFind.Prefix,
		Extension:   mediaToFind.Extension,
		ContentType: mediaToFind.ContentType,
		Url:         mediaToFind.Url,
		IsConfirm:   mediaToFind.IsConfirm,
		Timestamp:   mediaToFind.Timestamp,
	}, nil

}

// todo now only isconfirm update
func (repo *Database) UpdateMediaTempConfirm(userId string, mediaTempUpdate *mediaaggregate.MediaTempUpdate) error {
	prefix := mediaTempUpdate.MediaResource.Prefix.String()
	mediaType := mediaTempUpdate.MediaResource.ContentType.ConvertToMediaType()

	pk := GenerateUserDataPk(userId)
	sk := GenerateUserBoardImageDataSk(prefix, mediaType)

	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk},
		"Sk": &types.AttributeValueMemberS{Value: sk},
	}

	update := expression.Set(expression.Name("IsConfirm"), expression.Value(mediaTempUpdate.MediaEntity.IsConfirm)).
		Set(expression.Name("Timestamp"), expression.Value(mediaTempUpdate.MediaEntity.Timestamp))
	expr, err := GenerateUpdateExpression(update)
	if err != nil {
		return err
	}

	return repo.UpdateDynamoDBItem(key, expr)
}

func (repo *Database) DeleteMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error {
	key := CreateDeleteMediaData(userId, &mediaTemp.MediaResource)
	return repo.DeleteDynamoDBItem(key)
}

func CreateDeleteMediaData(userId string, mediaResource *mediavalueobject.MediaResource) map[string]types.AttributeValue {
	prefix := mediaResource.Prefix.String()
	mediaType := mediaResource.GetMediaType()

	pk := GenerateUserDataPk(userId)
	sk := GenerateUserBoardImageDataSk(prefix, mediaType)

	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk},
		"Sk": &types.AttributeValueMemberS{Value: sk},
	}

	return key
}
