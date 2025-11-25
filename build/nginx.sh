#!/bin/bash
set -e

# 前端服务构建脚本
# 用于 Zadig 构建配置，直接复制到"通用构建脚本"中

# 进入项目根目录（Zadig 的工作目录通常是 /workspace/，代码在 k8s_sample 子目录下）
cd k8s_sample

cd services/nginx
make build                    # 使用 Makefile 准备构建

# 回到项目根目录执行 Docker 构建
# Dockerfile 需要在项目根目录执行（因为需要访问 services/frontend 和 services/nginx）
cd ../..
docker build -t $IMAGE -f services/nginx/Dockerfile .
docker push $IMAGE

