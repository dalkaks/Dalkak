package database

import (
	"context"
	"dalkak/internal/infrastructure/database/dao"
	responseutil "dalkak/pkg/utils/response"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Database struct {
	client   *dynamodb.Client
	table    string
	queryKey string
}

const (
	UserIdEntityTypeIndex = "UserId-EntityType-index"
)

func NewDB(ctx context.Context, mode, queryKey string) (*Database, error) {
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

	return &Database{client: dbClient, table: table, queryKey: queryKey}, nil
}

func (db *Database) GetClient() *dynamodb.Client {
	return db.client
}

func (db *Database) GetTable() string {
	return db.table
}

func (db *Database) GetTransactionID() (*dao.TransactionDao, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"Pk": &types.AttributeValueMemberS{Value: "Server#Transaction"},
			"Sk": &types.AttributeValueMemberS{Value: "Server#Transaction"},
		},
	}

	result, err := db.client.GetItem(context.Background(), input)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	var transaction dao.TransactionDao
	err = attributevalue.UnmarshalMap(result.Item, &transaction)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}
	return &transaction, nil
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

func (db *Database) QueryItems(expr expression.Expression, index *string, pageDao *dao.RequestPageDao, dest interface{}) (*dao.ResponsePageDao, error) {
	var exclusiveStartKey map[string]types.AttributeValue

	if pageDao != nil && pageDao.ExclusiveStartKey != nil {
		decodedKey, err := decodeExclusiveStartKey(*pageDao.ExclusiveStartKey)
		if err != nil {
			return nil, err
		}
		exclusiveStartKey = decodedKey
	}

	var items []map[string]types.AttributeValue
	var count int

	for {
		input := &dynamodb.QueryInput{
			TableName:                 aws.String(db.table),
			IndexName:                 index,
			KeyConditionExpression:    expr.KeyCondition(),
			FilterExpression:          expr.Filter(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ExclusiveStartKey:         exclusiveStartKey,
		}

		if pageDao != nil {
			input.Limit = aws.Int32(int32(pageDao.Limit))
		}

		result, err := db.client.Query(context.Background(), input)
		if err != nil {
			return nil, err
		}

		items = append(items, result.Items...)
		count += len(result.Items)

		exclusiveStartKey = result.LastEvaluatedKey
		if pageDao != nil && len(result.Items) < pageDao.Limit && result.LastEvaluatedKey != nil {
			continue
		}

		break
	}

	err := attributevalue.UnmarshalListOfMaps(items, dest)
	if err != nil {
		return nil, err
	}

	var nextPageToken *string
	if len(exclusiveStartKey) > 0 {
		encodedKey, err := encodeExclusiveStartKey(exclusiveStartKey)
		if err != nil {
			return nil, err
		}
		nextPageToken = &encodedKey
	}

	return &dao.ResponsePageDao{
		Count:             count,
		ExclusiveStartKey: nextPageToken,
	}, nil
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

func (db *Database) WriteTransaction(builder *TransactionBuilder) error {
	_, err := db.client.TransactWriteItems(context.Background(), &dynamodb.TransactWriteItemsInput{
		TransactItems: builder.Build(),
	})
	if err != nil {
		var tce *types.TransactionCanceledException
		if errors.As(err, &tce) {
			return responseutil.NewAppError(responseutil.ErrCodeServiceDown, responseutil.ErrMsgDbInternalTrans, err)
		}
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	return nil
}

func encodeExclusiveStartKey(exclusiveStartKey map[string]types.AttributeValue) (string, error) {
	jsonBytes, err := json.Marshal(exclusiveStartKey)
	if err != nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	encodedString := base64.URLEncoding.EncodeToString(jsonBytes)
	return encodedString, nil
}

func decodeExclusiveStartKey(encodedKey string) (map[string]types.AttributeValue, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}

	var exclusiveStartKey map[string]types.AttributeValue
	err = json.Unmarshal(decodedBytes, &exclusiveStartKey)
	if err != nil {
		return nil, err
	}

	return exclusiveStartKey, nil
}
