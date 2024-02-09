package user

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/timeutils"

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
	newUser := &UserData{
		Pk:            GenerateUserDataPk(walletAddress),
		Sk:            GenerateUserDataPk(walletAddress),
		EntityType:    UserDataType,
		WalletAddress: walletAddress,
		Timestamp:     timeutils.GetTimestamp(),
	}

	av, err := attributevalue.MarshalMap(newUser)
	if err != nil {
		return err
	}

	_, err = repo.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(repo.table),
		Item:      av,
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepositoryImpl) FindUser(walletAddress string) (*dtos.UserDto, error) {
	var userToFind UserData
	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: GenerateUserDataPk(walletAddress)},
		"Sk": &types.AttributeValueMemberS{Value: GenerateUserDataPk(walletAddress)},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.table),
		Key:       key,
	}

	result, err := repo.client.GetItem(context.Background(), input)
	if err != nil {
		return nil, err
	}

	if result.Item != nil {
		err = attributevalue.UnmarshalMap(result.Item, &userToFind)
		if err != nil {
			return nil, err
		}
		return userToFind.ToUserDto(), nil
	}
	return nil, nil
}
