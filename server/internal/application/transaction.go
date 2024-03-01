package application

import (
	"dalkak/internal/infrastructure/database/dao"
	responseutil "dalkak/pkg/utils/response"
	"time"
)

type RetryableFunc[T any] func(txId string) (T, error)

func ExecuteOptimisticTransactionWithRetry[T any](app *ApplicationImpl, fn RetryableFunc[T]) (result T, err error) {
	const maxRetry = 3

	for attempt := 1; attempt <= maxRetry; attempt++ {
		var transacionItem *dao.TransactionDao
		transacionItem, err = app.Database.GetTransactionID()
		if err != nil {
			return
		}

		result, err = fn(transacionItem.Id)
		if err == nil {
			return
		}

		if !isTransactionError(err) {
			return
		}

		if attempt < maxRetry {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	return
}

func isTransactionError(err error) bool {
	if appError, ok := err.(*responseutil.AppError); ok {
		if appError.Code == responseutil.ErrCodeServiceDown {
			return true
		}
	}
	return false
}
