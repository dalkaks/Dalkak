package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	sdkConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS SDK config, %v", err)
		return err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key
		headObjectOutput, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Printf("Failed to get object head for bucket: %s, key: %s, %v", bucket, key, err)
			return nil
		}

		log.Printf("Received event for bucket: %s, key: %s, size: %d", bucket, key, headObjectOutput.ContentLength)
	}

	return nil
}

// func getUserMediaData(bucketName, objectKey string) (*user.UserMediaData, error) {
// 	// DynamoDB에서 데이터 조회 로직 구현
// 	// 예시 코드는 실제 로직과 다를 수 있음
// 	return nil, nil
// }

// func updateUserMediaData(data *user.UserMediaData) error {
// 	// DynamoDB에서 데이터 업데이트 로직 구현
// 	return nil
// }

// func deleteS3Object(bucketName, objectKey string) error {
// 	// S3 객체 삭제 로직 구현
// 	return nil
// }
