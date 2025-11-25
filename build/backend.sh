#!/bin/bash
set -e

# 后端服务构建脚本
# 用于 Zadig 构建配置，直接复制到"通用构建脚本"中

# 进入项目根目录（Zadig 的工作目录通常是 /workspace/，代码在 k8s_sample 子目录下）
cd k8s_sample

cd services/backend
make build                    # 使用 Makefile 构建应用

docker build -t $IMAGE -f Dockerfile .
docker push $IMAGE

