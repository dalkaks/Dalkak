#!/bin/bash

# 람다 파일명을 인자로 받아 확인
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <lambda_filename>"
    exit 1
fi

LAMBDA_NAME=$1

# 빌드 디렉토리로 이동
cd server/cmd/lambda/${LAMBDA_NAME} || exit

# Go 파일 빌드
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

# bootstrap을 zip 파일로 압축
zip bootstrap.zip bootstrap
if [ $? -ne 0 ]; then
    echo "Zip creation failed"
    exit 1
fi

echo "bootstrap.zip has been created successfully."

# 빌드 디렉토리로부터 원래 디렉토리로 돌아감
cd - > /dev/null

# ./scripts/build-lambda.sh s3tempconfirm