#!/bin/bash

# AWS CLI를 사용하여 ECS 서비스의 Desired Count를 0으로 설정하는 스크립트

# ECS 클러스터 이름
CLUSTER_NAME="dalkak-dev-cluster"
# ECS 서비스 이름
SERVICE_NAME="dalkak-dev-service"

# Desired Count를 0으로 설정하여 서비스를 중지
aws ecs update-service --cluster $CLUSTER_NAME --service $SERVICE_NAME --desired-count 0
