#!/bin/bash
set -e

# 后端服务构建脚本
# 用于 Zadig 构建配置，直接复制到"通用构建脚本"中

cd services/backend
make build                    # 使用 Makefile 构建应用

docker build -t $IMAGE -f Dockerfile .
docker push $IMAGE

