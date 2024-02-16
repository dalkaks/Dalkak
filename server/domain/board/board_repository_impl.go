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
