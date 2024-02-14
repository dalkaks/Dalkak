package user

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/dynamodbutils"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/timeutils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepositoryImpl struct {
	client *dynamodb.Client
	table  string
}

func NewUserRepository(db interfaces.Database) *UserRepositoryImpl {
	client := db.GetClient()
	table := db.GetTable()

	return &UserRepositoryImpl{
		client: client,
		table:  table,
	}
}

func (repo *UserRepositoryImpl) CreateUser(walletAddress string) error {
	pk := GenerateUserDataPk(walletAddress)
	newUser := &UserData{
		Pk:         pk,
		Sk:         pk,
		EntityType: UserDataType,
		Timestamp:  timeutils.GetTimestamp(),

		WalletAddress: walletAddress,
	}

	err := dynamodbutils.PutDynamoDBItem(repo.client, repo.table, newUser)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepositoryImpl) FindUser(walletAddress string) (*dtos.UserDto, error) {
	pk := GenerateUserDataPk(walletAddress)
	var userToFind *UserData

	keyCond := expression.Key("Pk").Equal(expression.Value(pk)).
		And(expression.Key("Sk").Equal(expression.Value(pk)))
	expr, err := dynamodbutils.GenerateQueryExpression(keyCond, nil)
	if err != nil {
		return nil, err
	}

	err = dynamodbutils.QuerySingleItem(repo.client, repo.table, expr, &userToFind)
	if err != nil || userToFind == nil {
		return nil, err
	}

	return userToFind.ToUserDto(), nil
}

func (repo *UserRepositoryImpl) CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error {
	mediaType, err := httputils.ConvertContentTypeToMediaType(dto.ContentType)
	if err != nil {
		return err
	}
	sk := GenerateUserBoardImageDataSk(dto.Prefix, mediaType)

	newUploadMedia := &UserMediaData{
		Pk:         GenerateUserDataPk(userId),
		Sk:         sk,
		EntityType: sk,
		Timestamp:  timeutils.GetTimestamp(),

		Id:          dto.ID,
		Prefix:      dto.Prefix,
		Extension:   dto.Extension,
		ContentType: dto.ContentType,
		Url:         dto.URL,
		IsConfirm:   false,
	}

	err = dynamodbutils.PutDynamoDBItem(repo.client, repo.table, newUploadMedia)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepositoryImpl) FindUserUploadMedia(userId string, dto *dtos.FindUserUploadMediaDto) (*dtos.MediaMeta, error) {
	sk := GenerateUserBoardImageDataSk(dto.Prefix, dto.MediaType.String())
	var mediaToFind *UserMediaData

	keyCond := expression.Key("Pk").Equal(expression.Value(GenerateUserDataPk(userId))).
		And(expression.Key("Sk").Equal(expression.Value(sk)))
	expr, err := dynamodbutils.GenerateQueryExpression(keyCond, nil)
	if err != nil {
		return nil, err
	}

	if dto.IsConfirm != nil && *dto.IsConfirm {
		filt := expression.Name("IsConfirm").Equal(expression.Value(true))
		expr, err = dynamodbutils.GenerateQueryExpression(keyCond, &filt)
	}
	if err != nil {
		return nil, err
	}

	err = dynamodbutils.QuerySingleItem(repo.client, repo.table, expr, &mediaToFind)
	if err != nil || mediaToFind == nil {
		return nil, err
	}

	return mediaToFind.ToMediaMeta(), nil
}

func (repo *UserRepositoryImpl) UpdateUserUploadMedia(userId string, findDto *dtos.MediaMeta, updateDto *dtos.UpdateUserUploadMediaDto) error {
	mediaType, err := httputils.ConvertContentTypeToMediaType(findDto.ContentType)
	if err != nil {
		return err
	}

	pk := GenerateUserDataPk(userId)
	sk := GenerateUserBoardImageDataSk(findDto.Prefix, mediaType)

	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: pk},
		"Sk": &types.AttributeValueMemberS{Value: sk},
	}

	update := expression.Set(expression.Name("IsConfirm"), expression.Value(updateDto.IsConfirm))
	expr, err := dynamodbutils.GenerateUpdateExpression(update)
	if err != nil {
		return err
	}

	return dynamodbutils.UpdateDynamoDBItem(repo.client, repo.table, key, expr)
}
