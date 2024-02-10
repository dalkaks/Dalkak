package board

import (
	"dalkak/pkg/interfaces"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type BoardRepositoryImpl struct {
	client *dynamodb.Client
	table  string
}

func NewBoardRepository(db interfaces.Database) *BoardRepositoryImpl {
	client := db.GetClient()
	table := db.GetTable()

	return &BoardRepositoryImpl{
		client: client,
		table:  table,
	}
}

// func (repo *BoardRepositoryImpl) CreateBoardImage(dto *dtos.BoardImageDto, boardId string) error {
// 	newBoardImage := &BoardImageData{
// 		Pk:          GenerateBoardDataPk(boardId),
// 		Sk:          BoardImageDataType + `#` + dto.Id,
// 		EntityType:  BoardImageDataType,
// 		Timestamp:   timeutils.GetTimestamp(),
// 		Id:          dto.Id,
// 		Extension:   dto.Extension,
// 		ContentType: dto.ContentType,
// 		Url:         dto.Url,
// 	}

// 	av, err := attributevalue.MarshalMap(newBoardImage)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = repo.client.PutItem(context.Background(), &dynamodb.PutItemInput{
// 		TableName: aws.String(repo.table),
// 		Item:      av,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
