# Microservice Demo (Go + UmiJS + Nginx + Kubernetes)

本项目从 0 到 1 构建了一个最小可用的微服务示例，包括：

- **Go 后端**：提供登录、注册、退出、用户列表接口，纯内存存储，启动时内置 `admin/admin` 并随机生成多名示例用户（每次启动都会生成不同的用户名）。
- **UmiJS 前端**：登录成功后展示带有“用户列表”侧边栏和顶部用户导航（点击用户名可退出），用户列表以表格形式呈现。
- **Nginx 入口**：负责托管前端静态资源并将 `/api` 流量反向代理到后端。
- **Kubernetes 清单**：描述三层架构（Namespace、Go 服务、Nginx Ingress）的部署方式。

## 目录结构

```
backend/                       # Go 后端（main.go, go.mod, Dockerfile）
frontend/                      # UmiJS 前端应用
nginx/                         # Nginx 配置与发布 Dockerfile
deploy/k8s/                    # Kubernetes namespace/deployment/service/ingress
```

## 本地开发

### 后端（Go）

```bash
cd backend
go run main.go
```

服务会监听 `:8080`，所有数据存储在内存中，并在启动时创建 `admin/admin` 以及若干随机示例用户。暴露的接口如下：

- `POST /api/register` body `{ "username": "...", "password": "..." }`
- `POST /api/login` -> `{ "token": "...", "message": "..." }`
- `POST /api/logout` 需 `Authorization: Bearer <token>`
- `GET /api/users` 需 `Authorization: Bearer <token>`，返回 `{ "users": [...], "me": "..." }`，结果中包含随机生成的示例用户和 `admin`

### 前端（UmiJS）

```bash
cd frontend
npm install
npm run dev
```

开发服务器默认会通过 `/api` 代理到 `http://localhost:8080`，因此确保 Go 服务已运行。

## 容器镜像

构建 Go 服务镜像：

```bash
docker build -t go-backend:latest backend
```

构建 Nginx（包含前端静态文件）镜像：

```bash
docker build -t micro-frontend-entry:latest -f nginx/Dockerfile .
```

`nginx/Dockerfile` 会在第一阶段安装前端依赖并执行 `npm run build`，然后将 `dist` 内容复制到轻量级 Nginx 镜像中，同时带入 `nginx/default.conf` 以将 `/api` 代理至 `go-backend:8080`。

## Kubernetes 部署

在打好镜像并推送到可访问的镜像仓库后，更新 `deploy/k8s/backend.yaml` 与 `deploy/k8s/nginx.yaml` 中的 `image` 字段，然后依次执行：

```bash
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/backend.yaml
kubectl apply -f deploy/k8s/nginx.yaml
```

资源说明：

- `Namespace micro-demo`：隔离微服务相关资源。
- `Deployment go-backend + Service go-backend`：暴露 8080 端口供 Nginx 访问。
- `Deployment web-entry + Service web-entry`：运行内置前端资源的 Nginx，监听 80。
- `Ingress web-entry`：示例中绑定 `micro.demo.local`。在本地集群中可通过 `hosts` 文件将该域名指向 Ingress Controller 的地址。

应用部署完成后，浏览器访问 `http://micro.demo.local/` 即可看到登录/注册界面，首次登录使用 `admin/admin`。

## 后续扩展

- 可将用户及 token 存储替换为 Redis、PostgreSQL 等外部服务。
- 为前端添加路由/状态管理或换成 Ant Design 等 UI 组件库。
- 在 Kubernetes 中增加 CI/CD、HPA、PodMonitor 等资源，实现真正的生产级架构。
