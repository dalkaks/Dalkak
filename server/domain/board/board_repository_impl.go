package board

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/timeutils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type BoardRepositoryImpl struct {
	client *dynamodb.Client
	prefix string
}

func NewBoardRepository(db interfaces.Database) *BoardRepositoryImpl {
	client := db.GetClient()
	prefix := db.GetPrefix()

	return &BoardRepositoryImpl{
		client: client,
		prefix: prefix,
	}
}

func (repo *BoardRepositoryImpl) CreateBoardImage(dto *dtos.BoardImageDto, userId string) error {
	table := repo.prefix + BoardImageTableName
	newBoardImage := &BoardImageTable{
		Id:          dto.Id,
		BoardId:     dto.BoardId,
		Extension:   dto.Extension,
		ContentType: dto.ContentType,
		Url:         dto.Url,
		UserId:      userId,
		Timestamp:   timeutils.GetTimestamp(),
	}

	av, err := attributevalue.MarshalMap(newBoardImage)
	if err != nil {
		return err
	}

	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	})
	if err != nil {
		return err
	}

	return nil
}
