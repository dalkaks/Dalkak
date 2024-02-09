package interfaces

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type Database interface {
	GetClient() *dynamodb.Client
	GetTable() string
}
