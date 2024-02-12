package dynamodbutils

import (
	"context"
	"dalkak/pkg/dtos"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func PutDynamoDBItem[T any](client *dynamodb.Client, tableName string, data T) error {
	av, err := attributevalue.MarshalMap(data)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to marshal data to map",
		}
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to put db data",
		}
	}

	return nil
}
