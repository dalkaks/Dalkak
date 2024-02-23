#!/bin/bash

SCHEMA_JSON_PATH="./config/mock/schema.json"

! lsof -i:8000 -t >/dev/null 2>&1 || (echo "Error: Port 8000 is already in use." && exit 1)

aws dynamodb describe-table --table-name dalkak_dev --region ap-northeast-2 --profile dalkak \
| jq '.Table | {
    AttributeDefinitions, 
    KeySchema, 
    TableName, 
    LocalSecondaryIndexes: (if .LocalSecondaryIndexes then .LocalSecondaryIndexes | map({
      IndexName,
      KeySchema,
      Projection
    }) else [] end),
    GlobalSecondaryIndexes: (if .GlobalSecondaryIndexes then .GlobalSecondaryIndexes | map({
      IndexName,
      KeySchema,
      Projection,
      ProvisionedThroughput: {
        ReadCapacityUnits: 1,
        WriteCapacityUnits: 1
      }
    }) else [] end)
  } + {
    ProvisionedThroughput: {
      ReadCapacityUnits: 1, 
      WriteCapacityUnits: 1
    }
  }' > ./config/mock/schema.json
  
docker run --name dynamodb-local --rm -d -p 8000:8000 amazon/dynamodb-local

echo "Waiting for DynamoDB Local to be ready..."
while ! timeout 10 bash -c "echo > /dev/tcp/localhost/8000"; do
  sleep 0.2
done
echo "DynamoDB Local is ready."

aws dynamodb create-table --cli-input-json file://${SCHEMA_JSON_PATH} --endpoint-url http://localhost:8000 --region ap-northeast-2 --profile dalkak > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "Table creation successful"
else
    echo "Table creation failed"
fi
