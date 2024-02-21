package database

import (
	userobject "dalkak/internal/domain/user/object"
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


func (repo *Database) CreateUser(user *userobject.UserEntity) error {
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

func (repo *Database) FindUserByWalletAddress(walletAddress string) (*userobject.UserEntity, error) {
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

	return &userobject.UserEntity{
		WalletAddress: userToFind.WalletAddress,
		Timestamp:     userToFind.Timestamp,
	}, nil
}

// func (repo *Database) CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error {
// 	mediaType, err := parseutil.ConvertContentTypeToMediaType(dto.ContentType)
// 	if err != nil {
// 		return err
// 	}
// 	sk := GenerateUserBoardImageDataSk(dto.Prefix, mediaType)

// 	newUploadMedia := &UserMediaData{
// 		Pk:         GenerateUserDataPk(userId),
// 		Sk:         sk,
// 		EntityType: sk,
// 		Timestamp:  timeutil.GetTimestamp(),

// 		Id:          dto.ID,
// 		Prefix:      dto.Prefix,
// 		Extension:   dto.Extension,
// 		ContentType: dto.ContentType,
// 		Url:         dto.URL,
// 		IsConfirm:   false,
// 	}

// 	err = PutDynamoDBItem(repo.client, repo.table, newUploadMedia)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (repo *Database) FindUserUploadMedia(userId string, dto *dtos.FindUserUploadMediaDto) (*dtos.MediaMeta, error) {
// 	sk := GenerateUserBoardImageDataSk(dto.Prefix, dto.MediaType.String())
// 	var mediaToFind *UserMediaData

// 	keyCond := expression.Key("Pk").Equal(expression.Value(GenerateUserDataPk(userId))).
// 		And(expression.Key("Sk").Equal(expression.Value(sk)))
// 	expr, err := GenerateQueryExpression(keyCond, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if dto.IsConfirm != nil && *dto.IsConfirm {
// 		filt := expression.Name("IsConfirm").Equal(expression.Value(true))
// 		expr, err = GenerateQueryExpression(keyCond, &filt)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = QuerySingleItem(repo.client, repo.table, expr, &mediaToFind)
// 	if err != nil || mediaToFind == nil {
// 		return nil, err
// 	}

// 	return mediaToFind.ToMediaMeta(), nil
// }

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
