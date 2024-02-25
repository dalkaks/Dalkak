package database

import (
	"context"
	responseutil "dalkak/pkg/utils/response"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Database struct {
	client *dynamodb.Client
	table  string
}

func NewDB(ctx context.Context, mode string) (*Database, error) {
	var cfg aws.Config
	var err error

	if mode == "LOCAL" {
		cfg, err = awsConfig.LoadDefaultConfig(ctx, awsConfig.WithSharedConfigProfile("dalkak"))
	} else {
		cfg, err = awsConfig.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	var table string
	if mode == "PROD" {
		table = "dalkak_prod"
	} else {
		table = "dalkak_dev"
	}

	return &Database{client: dbClient, table: table}, nil
}

func (db *Database) GetClient() *dynamodb.Client {
	return db.client
}

func (db *Database) GetTable() string {
	return db.table
}

func (db *Database) QuerySingleItem(expr expression.Expression, dest interface{}) error {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(db.table),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := db.client.Query(context.Background(), input)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	if len(result.Items) == 0 {
		return nil
	}

	err = attributevalue.UnmarshalMap(result.Items[0], dest)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}
	return nil
}

func (db *Database) PutDynamoDBItem(data interface{}) error {
	av, err := attributevalue.MarshalMap(data)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	createExpr, err := GenerateCreateExpression()
	if err != nil {
		return err
	}

	_, err = db.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:                 aws.String(db.table),
		Item:                      av,
		ExpressionAttributeNames:  createExpr.Names(),
		ExpressionAttributeValues: createExpr.Values(),
		ConditionExpression:       createExpr.Condition(),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			return responseutil.NewAppError(responseutil.ErrCodeConflict, responseutil.ErrMsgDataConflict, err)
		}
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return nil
}

func (db *Database) UpdateDynamoDBItem(key map[string]types.AttributeValue, expr expression.Expression) error {
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.table),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueNone,
	}

	_, err := db.client.UpdateItem(context.Background(), input)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return nil
}

func (db *Database) DeleteDynamoDBItem(key map[string]types.AttributeValue) error {
	_, err := db.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(db.table),
		Key:       key,
	})
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return nil
}

func GenerateQueryExpression(keyCond expression.KeyConditionBuilder, filt *expression.ConditionBuilder) (expression.Expression, error) {
	builder := expression.NewBuilder().WithKeyCondition(keyCond)

	if filt != nil {
		builder = builder.WithFilter(*filt)
	}

	expr, err := builder.Build()
	if err != nil {
		return expression.Expression{}, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return expr, nil
}

func GenerateCreateExpression() (expression.Expression, error) {
	condition := expression.AttributeNotExists(expression.Name("Pk")).And(expression.AttributeNotExists(expression.Name("Sk")))
	builder := expression.NewBuilder().WithCondition(condition)

	expr, err := builder.Build()
	if err != nil {
		return expression.Expression{}, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return expr, nil
}

func GenerateUpdateExpression(update expression.UpdateBuilder) (expression.Expression, error) {
	condition := expression.AttributeExists(expression.Name("Pk")).And(expression.AttributeExists(expression.Name("Sk")))
	builder := expression.NewBuilder().WithUpdate(update).WithCondition(condition)

	expr, err := builder.Build()
	if err != nil {
		return expression.Expression{}, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return expr, nil
}
