#!/bin/bash

SCHEMA_JSON_PATH="./config/mock/schema.json"
DATABASE_JSON_PATH="./config/mock/database.json"

! lsof -i:8000 -t >/dev/null 2>&1 || (echo "Error: Port 8000 is already in use." && exit 1)

aws dynamodb describe-table --table-name dalkak_dev --region ap-northeast-2 --output json > ./config/mock/schema.json --profile dalkak
aws dynamodb scan --table-name dalkak_dev --region ap-northeast-2 \
| jq '{"dalkak-dev": [.Items[] |  {PutRequest: {Item: .}}]}' > ./config/mock/schema.json --profile dalkak

docker run --name dynamodb-local --rm -d -p 8000:8000 amazon/dynamodb-local

aws dynamodb create-table --cli-input-json file://${SCHEMA_JSON_PATH} --endpoint-url http://localhost:8000 --region ap-northeast-2 --profile dalkak
aws dynamodb batch-write-item --request-items file://${DATABASE_JSON_PATH} --endpoint-url http://localhost:8000 --region ap-northeast-2 --profile dalkak

echo "DynamoDB Local has been started successfully."