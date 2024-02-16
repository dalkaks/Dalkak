package dynamodbutils

import (
	"context"
	"dalkak/pkg/dtos"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutDynamoDBItem[T any](client *dynamodb.Client, tableName string, data T) error {
	av, err := attributevalue.MarshalMap(data)
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgDBInternal, err)
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgDBInternal, err)
	}

	return nil
}

func UpdateDynamoDBItem(client *dynamodb.Client, tableName string, key map[string]types.AttributeValue, expr expression.Expression) error {
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueNone,
	}

	_, err := client.UpdateItem(context.Background(), input)
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgDBInternal, err)
	}

	return nil
}

func DeleteDynamoDBItem(client *dynamodb.Client, tableName string, key map[string]types.AttributeValue) error {
	_, err := client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgDBInternal, err)
	}

	return nil
}
