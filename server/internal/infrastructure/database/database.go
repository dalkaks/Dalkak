package database

import (
	"context"
	"dalkak/internal/infrastructure/database/dao"
	cryptoutil "dalkak/pkg/utils/crypto"
	responseutil "dalkak/pkg/utils/response"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

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
	queryKey []byte
}

const (
	UserIdEntityTypeIndex    = "UserId-EntityType-index"
	EntityTypeTimestampIndex = "EntityType-Timestamp-index"
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

	queryKeyByte := []byte(queryKey)

	return &Database{client: dbClient, table: table, queryKey: queryKeyByte}, nil
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
	var err error

	if pageDao != nil && pageDao.ExclusiveStartKey != nil && *pageDao.ExclusiveStartKey != "" {
		exclusiveStartKey, err = db.decryptExclusiveStartKey(*pageDao.ExclusiveStartKey)
		if err != nil {
			return nil, err
		}
	}

	var items []map[string]types.AttributeValue
	var count int
	var nextPageToken *string

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

		if len(result.LastEvaluatedKey) > 0 {
			exclusiveStartKey = result.LastEvaluatedKey
			encryptedKey, err := db.encryptExclusiveStartKey(exclusiveStartKey)
			if err != nil {
				return nil, err
			}
			nextPageToken = &encryptedKey
		} else {
			nextPageToken = nil
			break
		}

		if pageDao != nil && pageDao.Limit > 0 && count >= pageDao.Limit {
			break
		}
	}

	err = attributevalue.UnmarshalListOfMaps(items, dest)
	if err != nil {
		return nil, err
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
		TableName:                           aws.String(db.table),
		Key:                                 key,
		UpdateExpression:                    expr.Update(),
		ExpressionAttributeNames:            expr.Names(),
		ExpressionAttributeValues:           expr.Values(),
		ConditionExpression:                 expr.Condition(),
		ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureNone,
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

func (db *Database) encryptExclusiveStartKey(exclusiveStartKey map[string]types.AttributeValue) (string, error) {
	byteKey, err := json.Marshal(exclusiveStartKey)
	if err != nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	encryptKey, err := cryptoutil.EncryptAES(db.queryKey, byteKey)
	if err != nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	encodedString := base64.URLEncoding.EncodeToString(encryptKey)
	return encodedString, nil
}

func (db *Database) decryptExclusiveStartKey(encodedKey string) (map[string]types.AttributeValue, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}

	encryptKey := []byte(decodedBytes)
	decryptKey, err := cryptoutil.DecryptAES(db.queryKey, encryptKey)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	var tempMap map[string]map[string]interface{}
	err = json.Unmarshal(decryptKey, &tempMap)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	decodedKey := make(map[string]types.AttributeValue)
	for k, v := range tempMap {
		if num, err := strconv.ParseFloat(v["Value"].(string), 64); err == nil {
			decodedKey[k] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(num, 'f', -1, 64)}
		} else {
			decodedKey[k] = &types.AttributeValueMemberS{Value: v["Value"].(string)}
		}
	}
	return decodedKey, nil
}
