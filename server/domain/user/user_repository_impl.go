package user

import (
	"context"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/timeutils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepositoryImpl struct {
	client *dynamodb.Client
	prefix string
}

func NewUserRepository(db interfaces.Database) *UserRepositoryImpl {
	client := db.GetClient()
	prefix := db.GetPrefix()

	return &UserRepositoryImpl{
		client: client,
		prefix: prefix,
	}
}

func (repo *UserRepositoryImpl) FindOrCreateUser(walletAddress string) (string, error) {
	table := repo.prefix + UserTableName
	User := UserTable{WalletAddress: walletAddress}

	response, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			WalletAddressKey: &types.AttributeValueMemberS{Value: User.WalletAddress},
		},
	})
	if err != nil {
		return "", err
	}

	if response.Item != nil {
		err = attributevalue.UnmarshalMap(response.Item, &User)
		if err != nil {
			return "", err
		}
		return User.WalletAddress, nil
	}

	newUser := UserTable{
		WalletAddress: walletAddress,
		Timestamp:     timeutils.GetTimestamp(),
	}

	av, err := attributevalue.MarshalMap(newUser)
	if err != nil {
		return "", err
	}

	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	})
	if err != nil {
		return "", err
	}

	return newUser.WalletAddress, nil
}
