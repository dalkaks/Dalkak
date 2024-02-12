package dynamodbutils

import (
	"context"
	"dalkak/pkg/dtos"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func QuerySingleItem(client *dynamodb.Client, tableName string, expr expression.Expression, dest interface{}) error {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := client.Query(context.Background(), input)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to query user media data",
		}
	}

	if len(result.Items) == 0 {
		return nil
	}

	err = attributevalue.UnmarshalMap(result.Items[0], dest)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to unmarshal user media data",
		}
	}
	return nil
}
