# 使用 Zadig 部署微服务 Demo

本文描述如何通过 [Zadig](https://zadig.fit2cloud.com/) 将本项目（Go 后端 + UmiJS 前端 + Nginx 入口）部署到 Kubernetes 集群。假设你已经具备：

- 一个可访问的 Kubernetes 集群，以及能连接到该集群的 Zadig 环境。
- 已安装并配置好的容器镜像仓库（如 Harbor、ACR、ECR 等）。
- Zadig 平台上创建好了产品（Product）并绑定到目标集群和命名空间。

## 1. 准备代码仓库

整个项目位于根目录 `go_service_sample/`，主要模块如下：

```
backend/    # Go 服务 (main.go、go.mod、Dockerfile)
frontend/   # UmiJS 前端源代码
nginx/      # Nginx Dockerfile + default.conf（使用前端构建产物）
deploy/k8s/ # Kubernetes YAML（Namespace、后端 Deployment、入口 Deployment/Ingress）
```

将该目录作为一个 Git 仓库连接到 Zadig 的代码源即可。

## 2. 配置构建与镜像

在 Zadig 中分别为后端和前端创建两个「服务组件」：

### 后端服务（go-backend）

1. **源码目录**：`backend/`
2. **构建命令**（示例）：

   ```bash
   cd backend
   go mod tidy
   go build -o server .
   ```

3. **Dockerfile**：使用 `backend/Dockerfile`
4. **镜像名称**：`<REGISTRY>/micro-demo/go-backend:<tag>`
5. **容器端口**：8080

### 前端入口（micro-frontend-entry）

1. **源码目录**：`frontend/`
2. **构建命令**：

   ```bash
   cd frontend
   npm install
   npm run build
   ```

3. **Dockerfile**：使用 `nginx/Dockerfile`（第一阶段需要前端依赖，第二阶段拷贝至 Nginx）
4. **镜像名称**：`<REGISTRY>/micro-demo/frontend-entry:<tag>`
5. **容器端口**：80

> 注：`nginx/Dockerfile` 默认会将 `/api` 请求代理到名为 `go-backend:8080` 的服务。保持 Kubernetes Service 名称一致即可。

构建完成后，确保镜像成功推送到你的镜像仓库。

## 3. 在 Zadig 中定义环境

### 3.1 创建 Namespace

可直接使用仓库中的 `deploy/k8s/namespace.yaml` 在集群中创建 `micro-demo` 命名空间，也可以在 Zadig 平台上引用已存在的 Namespace。

### 3.2 部署服务

在 Zadig 产品中添加两个服务：

1. **go-backend**
   - Deployment 模板可参考 `deploy/k8s/backend.yaml`。
   - 更新模板中的 `image` 字段为上一步推送的镜像地址。
   - 保留 `Service` 配置（ClusterIP + 8080）以及 `readiness/liveness probe`。

2. **web-entry**
   - Deployment/Service/Ingress 模板可参考 `deploy/k8s/nginx.yaml`。
   - 修改 `image` 字段为前端镜像地址。
   - 如果集群已有 Ingress Controller，保持 `Ingress` 配置；否则可删除或替换为 LoadBalancer。

### 3.3 环境变量与配置

当前示例中所有数据存储于内存，无需额外配置。若未来接入数据库或外部配置，可在 Zadig 服务的「环境变量」或「ConfigMap/Secret」中补充。

## 4. 发布流程

1. 在 Zadig 的工作流中添加两个构建 Job（后端、前端）与一个部署 Job，串联顺序如下：
   1. 后端构建并推送镜像
   2. 前端构建并推送镜像
   3. 更新环境：go-backend → web-entry
2. 触发工作流后，等待镜像构建、部署完成。
3. 部署完成后，在浏览器中访问 `http://micro.demo.local/`（或你设置的 Ingress 域名），使用默认账号 `admin/admin` 登录即可查看带侧边栏 + 表格的后台页面。

## 5. 常见调整

- **自定义域名**：修改 `deploy/k8s/nginx.yaml` 中 Ingress 的 `host`，并在 DNS/hosts 中解析到 Ingress Controller。
- **用户初始数据**：后端会在启动时生成多个随机用户，若希望固定账号可修改 `backend/main.go` 的 `seedUsers` 逻辑。
- **横向扩容**：在 Zadig 服务配置中调整 Deployment 的 `replicas` 或启用 HPA。
- **日志与监控**：可在 Zadig 环境中接入 Prometheus/EFK 等工具来扩展观测能力。

## 6. 以“项目”形式部署

若你的 Zadig 已启用「项目(Project)」模式（通常在多模块协作或多环境共用场景），可以将该仓库作为单个项目接入，流程示例如下：

1. **创建项目**：在 Zadig 后台创建一个新项目（例如 `micro-demo`），选择代码仓库并指定默认分支。
2. **定义服务模板**：
   - 在项目中添加两个服务模板：`go-backend`（来源 `backend/`）与 `web-entry`（来源 `frontend/` + `nginx/`）。
   - 为每个模板配置构建命令、Dockerfile 路径及镜像仓库推送信息。
3. **创建环境**：
   - 在项目下创建测试/生产等环境，指定 Kubernetes 集群及 Namespace（可直接引用 `micro-demo`）。
   - 将服务模板添加到对应环境中，Zadig 会自动生成 Deployment/Service/Ingress。
4. **构建与部署**：
   - 在项目中创建工作流：步骤 1/2 分别构建后端、前端镜像；步骤 3 发布到目标环境。
   - 如果需要按环境区分镜像标签，可在工作流参数中传入 `TAG`，构建阶段引用 `TAG` 生成镜像，发布阶段统一使用。
5. **多环境推广**：
   - 可在项目中配置多条工作流，例如「Dev → Staging → Prod」推广链路，通过镜像复用或 Helm/Kustomize 参数差异化部署。

此模式下，所有服务、环境与流水线均在同一个项目下管理，便于团队协作与统一审计。

至此，整个项目即可通过 Zadig 完成自动构建、发布与迭代。如需更多高级玩法（多环境、灰度发布、工作流嵌套等），请参考 Zadig 官方文档。祝部署顺利! 🎉
