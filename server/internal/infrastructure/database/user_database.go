package database

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	userentity "dalkak/internal/domain/user/object/entity"
	"dalkak/internal/infrastructure/database/dao"
	responseutil "dalkak/pkg/utils/response"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

func (repo *Database) CreateUserMediaTemp(userId string, mediaTemp *mediaaggregate.MediaTempAggregate) error {
	sk := GenerateUserBoardImageDataSk(mediaTemp.Prefix.String(), mediaTemp.ContentType.ConvertToMediaType())

	newUploadMedia := &UserMediaData{
		Pk:         GenerateUserDataPk(userId),
		Sk:         sk,
		EntityType: sk,
		Timestamp:  mediaTemp.MediaEntity.Timestamp,

		Id:          mediaTemp.MediaEntity.Id,
		Prefix:      mediaTemp.Prefix.String(),
		Extension:   mediaTemp.ContentType.ConvertToExtension(),
		ContentType: mediaTemp.ContentType.String(),
		Url:         mediaTemp.MediaTempUrl.AccessUrl,
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

// func (repo *Database) UpdateUserUploadMedia(userId string, findDto *dtos.MediaMeta, updateDto *dtos.UpdateUserUploadMediaDto) error {
// 	mediaType, err := parseutil.ConvertContentTypeToMediaType(findDto.ContentType)
// 	if err != nil {
// 		return err
// 	}

// 	pk := GenerateUserDataPk(userId)
// 	sk := GenerateUserBoardImageDataSk(findDto.Prefix, mediaType)

// 	key := map[string]types.AttributeValue{
// 		"Pk": &types.AttributeValueMemberS{Value: pk},
// 		"Sk": &types.AttributeValueMemberS{Value: sk},
// 	}

// 	update := expression.Set(expression.Name("IsConfirm"), expression.Value(updateDto.IsConfirm))
// 	expr, err := GenerateUpdateExpression(update)
// 	if err != nil {
// 		return err
// 	}

// 	return UpdateDynamoDBItem(repo.client, repo.table, key, expr)
// }

// func (repo *Database) DeleteUserUploadMedia(userId string, dto *dtos.MediaMeta) error {
// 	mediaType, err := parseutil.ConvertContentTypeToMediaType(dto.ContentType)
// 	if err != nil {
// 		return err
// 	}

// 	pk := GenerateUserDataPk(userId)
// 	sk := GenerateUserBoardImageDataSk(dto.Prefix, mediaType)

// 	key := map[string]types.AttributeValue{
// 		"Pk": &types.AttributeValueMemberS{Value: pk},
// 		"Sk": &types.AttributeValueMemberS{Value: sk},
// 	}

// 	return DeleteDynamoDBItem(repo.client, repo.table, key)
// }
