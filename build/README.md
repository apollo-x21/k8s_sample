# 构建脚本说明

本目录包含用于 Zadig 构建配置的脚本文件。

## 使用方法

### 在 Zadig 中配置构建

1. 进入 Zadig 项目的服务配置页面
2. 点击"添加构建"或"编辑构建"
3. 在"通用构建脚本"文本框中，**直接复制对应的脚本内容**

### 脚本文件

- **backend.sh** - 后端服务构建脚本
  - 使用 Makefile 构建 Go 应用
  - 构建并推送 Docker 镜像

- **nginx.sh** - 前端服务构建脚本
  - 构建包含前端静态文件的 Nginx 镜像
  - 构建并推送 Docker 镜像

## 注意事项

1. **$IMAGE 变量**：这是 Zadig 自动提供的环境变量，包含完整的镜像地址（仓库+标签），无需修改

2. **路径**：脚本中的路径基于项目根目录，确保在 Zadig 中配置的代码库路径正确

3. **依赖**：
   - backend.sh 需要 Go 环境和 Makefile
   - nginx.sh 需要 Docker 环境

## 快速复制

### 后端构建脚本

```bash
#!/bin/bash
set -e

cd services/backend
make build

docker build -t $IMAGE -f Dockerfile .
docker push $IMAGE
```

### 前端构建脚本

```bash
#!/bin/bash
set -e

docker build -t $IMAGE -f services/nginx/Dockerfile .
docker push $IMAGE
```

