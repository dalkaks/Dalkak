package database

import (
	responseutil "dalkak/pkg/utils/response"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

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
