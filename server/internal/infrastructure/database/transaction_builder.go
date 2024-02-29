package database

import (
	generateutil "dalkak/pkg/utils/generate"
	responseutil "dalkak/pkg/utils/response"
	timeutil "dalkak/pkg/utils/time"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TransactionBuilder struct {
	TransactItems []types.TransactWriteItem
	tableName     string
}

func NewTransactionBuilder(tableName string, transactionId ...string) *TransactionBuilder {
	builder := TransactionBuilder{
		TransactItems: []types.TransactWriteItem{},
		tableName:     tableName,
	}

	if len(transactionId) > 0 {
		builder.AddTransaction(tableName, transactionId[0])
	}

	return &builder
}

func (builder *TransactionBuilder) Build() []types.TransactWriteItem {
	return builder.TransactItems
}

func (builder *TransactionBuilder) AddTransaction(tableName, exTxId string) {
	newTxId := generateutil.GenerateUUID()
	timestampStr := strconv.FormatInt(timeutil.GetTimestamp(), 10)

	key := map[string]types.AttributeValue{
		"Pk": &types.AttributeValueMemberS{Value: "Server#Transaction"},
		"Sk": &types.AttributeValueMemberS{Value: "Server#Transaction"},
	}

	updateExpression := "SET Id = :newTxId, #ts = :timestamp"
	conditionExpression := "Id = :exTxId"
	expressionAttributeNames := map[string]string{
		"#ts": "Timestamp",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
		":newTxId":   &types.AttributeValueMemberS{Value: newTxId},
		":exTxId":    &types.AttributeValueMemberS{Value: exTxId},
		":timestamp": &types.AttributeValueMemberN{Value: timestampStr},
	}

	builder.TransactItems = append(builder.TransactItems, types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 &tableName,
			Key:                       key,
			UpdateExpression:          &updateExpression,
			ConditionExpression:       &conditionExpression,
			ExpressionAttributeNames:  expressionAttributeNames,
			ExpressionAttributeValues: expressionAttributeValues,
		},
	})
}

func (builder *TransactionBuilder) AddPutItem(item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal, err)
	}

	createExpr, err := GenerateCreateExpression()
	if err != nil {
		return err
	}

	builder.TransactItems = append(builder.TransactItems, types.TransactWriteItem{
		Put: &types.Put{
			TableName:                 aws.String(builder.tableName),
			Item:                      av,
			ExpressionAttributeNames:  createExpr.Names(),
			ExpressionAttributeValues: createExpr.Values(),
			ConditionExpression:       createExpr.Condition(),
		},
	})
	return nil
}
