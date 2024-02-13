package main

import (
	"bytes"
	"context"
	"dalkak/config"
	"dalkak/pkg/utils/httputils"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

	lambdaConfig, err := config.LoadConfig[config.LambdaConfig](ctx, "Common", "LambdaConfig")
	devDomain := lambdaConfig.DevDomain
	prodDomain := lambdaConfig.ProdDomain

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
			continue
		}

		log.Printf("Received event for bucket: %s, key: %s, metadata: %v", bucket, key, headObjectOutput.Metadata)

		userId := headObjectOutput.Metadata["userid"]
		if userId == "" {
			log.Printf("No user id in metadata, bucket: %s, key: %s", bucket, key)
			continue
		}

		prefix := extractPrefix(key)
		if prefix == "" {
			log.Printf("No prefix in key, bucket: %s, key: %s", bucket, key)
			continue
		}

		if headObjectOutput.ContentType == nil {
			log.Printf("No content type in metadata, bucket: %s, key: %s", bucket, key)
			continue
		}
		mediaType, err := httputils.ConvertContentTypeToMediaType(*headObjectOutput.ContentType)
		if err != nil {
			log.Printf("Failed to convert content type to media type, %v", err)
			continue
		}

		var requestUrl string
		if bucket == "dalkak-dev" {
			requestUrl = devDomain + "/user/media/confirm"
		} else if bucket == "dalkak-prod" {
			requestUrl = prodDomain + "/user/media/confirm"
		} else {
			log.Printf("Invalid bucket name, %s", bucket)
			continue
		}

		requestBody, err := json.Marshal(map[string]string{
			"userId":    userId,
			"prefix":    prefix,
			"mediaType": mediaType,
		})
		if err != nil {
			log.Printf("Failed to marshal request body, %v", err)
			continue
		}

		resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Failed to send http request, %v", err)
			continue
		}
		defer resp.Body.Close()

		log.Printf("resp: %v", resp)
	}
	return nil
}

func extractPrefix(key string) string {
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
