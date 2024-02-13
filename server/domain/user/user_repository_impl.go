package user

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/dynamodbutils"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/timeutils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	Pk := GenerateUserDataPk(walletAddress)
	newUser := &UserData{
		Pk:         Pk,
		Sk:         Pk,
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
	Pk := GenerateUserDataPk(walletAddress)
	var userToFind *UserData

	keyCond := expression.Key("Pk").Equal(expression.Value(Pk)).
		And(expression.Key("Sk").Equal(expression.Value(Pk)))
	expr, err := dynamodbutils.GenerateExpression(keyCond, nil)
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
	Sk := GenerateUserBoardImageDataSk(dto.Prefix, mediaType)

	newUploadMedia := &UserMediaData{
		Pk:         GenerateUserDataPk(userId),
		Sk:         Sk,
		EntityType: Sk,
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
	Sk := GenerateUserBoardImageDataSk(dto.Prefix, dto.MediaType.String())
	var mediaToFind *UserMediaData

	keyCond := expression.Key("Pk").Equal(expression.Value(GenerateUserDataPk(userId))).
		And(expression.Key("Sk").Equal(expression.Value(Sk)))
	expr, err := dynamodbutils.GenerateExpression(keyCond, nil)
	if err != nil {
		return nil, err
	}

	if dto.IsConfirm != nil && *dto.IsConfirm {
		filt := expression.Name("IsConfirm").Equal(expression.Value(true))
		expr, err = dynamodbutils.GenerateExpression(keyCond, &filt)
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
