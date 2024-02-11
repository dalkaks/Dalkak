package user

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/timeutils"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	Pk := GenerateUserDataPk(walletAddress)
	newUser := &UserData{
		Pk:            Pk,
		Sk:            Pk,
		EntityType:    UserDataType,
		WalletAddress: walletAddress,
		Timestamp:     timeutils.GetTimestamp(),
	}

	av, err := attributevalue.MarshalMap(newUser)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to marshal user data to map",
		}
	}

	_, err = repo.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(repo.table),
		Item:      av,
	})
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to put user data",
		}
	}
	return nil
}

func (repo *UserRepositoryImpl) FindUser(walletAddress string) (*dtos.UserDto, error) {
	Pk := GenerateUserDataPk(walletAddress)
	var userToFind UserData
	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: Pk},
		"Sk": &types.AttributeValueMemberS{Value: Pk},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.table),
		Key:       key,
	}

	result, err := repo.client.GetItem(context.Background(), input)
	if err != nil {
		return nil, &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user data",
		}
	}

	if result.Item != nil {
		err = attributevalue.UnmarshalMap(result.Item, &userToFind)
		if err != nil {
			return nil, &dtos.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to unmarshal user data",
			}
		}
		return userToFind.ToUserDto(), nil
	}
	return nil, nil
}

func (repo *UserRepositoryImpl) CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error {
	mediaType, err := ConvertContentTypeToMediaType(dto.ContentType)
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

	av, err := attributevalue.MarshalMap(newUploadMedia)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to marshal user media data to map",
		}
	}

	_, err = repo.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(repo.table),
		Item:      av,
	})
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to put user media data",
		}
	}
	return nil
}

func (repo *UserRepositoryImpl) FindUserUploadMedia(userId string, dto *payloads.UserGetMediaRequest) (*dtos.MediaMeta, error) {
	Sk := GenerateUserBoardImageDataSk(dto.Prefix, dto.MediaType)
	var mediaToFind UserMediaData
	key := map[string]types.AttributeValue{
		"Pk":        &types.AttributeValueMemberS{Value: GenerateUserDataPk(userId)},
		"Sk":        &types.AttributeValueMemberS{Value: Sk},
		"IsConfirm": &types.AttributeValueMemberBOOL{Value: true},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.table),
		Key:       key,
	}

	result, err := repo.client.GetItem(context.Background(), input)
	if err != nil {
		return nil, &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user media data",
		}
	}

	if result.Item != nil {
		err = attributevalue.UnmarshalMap(result.Item, &mediaToFind)
		if err != nil {
			return nil, &dtos.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to unmarshal user media data",
			}
		}
		return mediaToFind.ToMediaMeta(), nil
	}
	return nil, nil
}
