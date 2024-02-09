#!/bin/bash

# AWS CLI를 사용하여 ECS 서비스의 Desired Count를 설정하는 스크립트

# ECS 클러스터 이름
CLUSTER_NAME="dalkak-dev-cluster"
# ECS 서비스 이름
SERVICE_NAME="dalkak-dev-service"

# 첫 번째 파라미터로 전달된 값을 사용
DESIRED_COUNT=${1:-0}

# 입력 값 검증 및 기본값 설정
if [[ $DESIRED_COUNT == !"0" || $DESIRED_COUNT == !"1" ]]; then
  echo "Invalid input. Only '0' or '1' is allowed. Defaulting to '0'."
  exit 1
fi

# Desired Count를 설정하여 서비스 업데이트
aws ecs update-service --cluster $CLUSTER_NAME --service $SERVICE_NAME --desired-count $DESIRED_COUNT --profile dalkak
