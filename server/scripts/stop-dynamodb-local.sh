#!/bin/bash

docker stop dynamodb-local || echo "DynamoDB Local container not found. Maybe it's already stopped."