package dynamodbutils

import (
	"dalkak/pkg/dtos"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func GenerateQueryExpression(keyCond expression.KeyConditionBuilder, filt *expression.ConditionBuilder) (expression.Expression, error) {
	builder := expression.NewBuilder().WithKeyCondition(keyCond)

	if filt != nil {
		builder = builder.WithFilter(*filt)
	}

	expr, err := builder.Build()
	if err != nil {
		return expression.Expression{}, &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to build db expression",
		}
	}

	return expr, nil
}

func GenerateUpdateExpression(update expression.UpdateBuilder) (expression.Expression, error) {
	builder := expression.NewBuilder().WithUpdate(update)

	expr, err := builder.Build()
	if err != nil {
		return expression.Expression{}, &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to build db expression",
		}
	}

	return expr, nil
}
