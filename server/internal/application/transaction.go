package application

import (
	responseutil "dalkak/pkg/utils/response"
	"time"
)

type RetryableFunc[T any] func(txId string) (T, error)

func ExecuteTransaction[T any](app *ApplicationImpl, fn RetryableFunc[T]) (T, error) {
	const maxRetry = 3

	for attempt := 1; attempt <= maxRetry; attempt++ {
		transacionItem, err := app.Database.GetTransactionID()
		if err != nil {
			return *new(T), err
		}

		result, err := fn(transacionItem.Id)
		if err == nil {
			return result, nil
		}

		if !isTransactionError(err) {
			return *new(T), err
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	// todo log and alarm
	return *new(T), responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgDBInternal)
}

func isTransactionError(err error) bool {
	// todo is transaction error
	return false
}
