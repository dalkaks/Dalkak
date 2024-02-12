package dynamodbutils

import (
	"dalkak/pkg/dtos"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func GenerateExpression(keyCond expression.KeyConditionBuilder, filt *expression.ConditionBuilder) (expression.Expression, error) {
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
