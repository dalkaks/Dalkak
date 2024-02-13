package main

import (
	"bytes"
	"context"
	"dalkak/config"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/lambdautils"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	sdkConfig, lambdaConfig, err := lambdautils.GetLambdaSdkConfig(ctx)
	if err != nil {
		lambdautils.HandleError(err)
		return err
	}

	s3Client := s3.NewFromConfig(*sdkConfig)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key
		headObjectOutput, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			lambdautils.HandleError(err)
			continue
		}
		log.Printf("Received event for bucket: %s, key: %s, metadata: %v", bucket, key, headObjectOutput.Metadata)

		userId := headObjectOutput.Metadata["userid"]
		if userId == "" {
			lambdautils.HandleError(errors.New("No user id in metadata, bucket: " + bucket + ", key: " + key))
			continue
		}

		if headObjectOutput.ContentType == nil {
			lambdautils.HandleError(errors.New("No content type in head object output, bucket: " + bucket + ", key: " + key))
			continue
		}
		mediaType, err := httputils.ConvertContentTypeToMediaType(*headObjectOutput.ContentType)
		if err != nil {
			lambdautils.HandleError(errors.New("Failed to convert content type to media type, bucket: " + bucket + ", key: " + key + ", content type: " + *headObjectOutput.ContentType))
			continue
		}

		requestUrl, err := getRequestUrl(bucket, lambdaConfig)
		if err != nil {
			lambdautils.HandleError(err)
			continue
		}
		requestBody, err := getRequestBody(userId, key, mediaType)
		if err != nil {
			lambdautils.HandleError(err)
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

func getRequestUrl(bucket string, lambdaConfig *config.LambdaConfig) (string, error) {
	if bucket == "dalkak-dev" {
		return lambdaConfig.DevDomain + "/user/media/confirm", nil
	} else if bucket == "dalkak-prod" {
		return lambdaConfig.ProdDomain + "/user/media/confirm", nil
	} else {
		return "", errors.New("Invalid bucket name, " + bucket)
	}
}

func getRequestBody(userId, key, mediaType string) ([]byte, error) {
	requestData := payloads.UserConfirmMediaRequest{
		UserId:    userId,
		Key:       key,
		MediaType: mediaType,
	}
	return json.Marshal(requestData)
}
