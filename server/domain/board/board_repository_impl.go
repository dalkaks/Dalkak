package board

import (
	"dalkak/pkg/interfaces"

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
