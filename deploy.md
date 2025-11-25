# 使用 Zadig 部署微服务 Demo

本文描述如何通过 [Zadig](https://zadig.fit2cloud.com/) 将本项目（Go 后端 + UmiJS 前端 + Nginx 入口）部署到 Kubernetes 集群。

参考文档：[Zadig 官方教程](https://koderover.com/tutorials-detail/codelabs/t100/index.html)

## 前置条件

在开始之前，请确保你已经具备：

- 一个可访问的 Kubernetes 集群，以及能连接到该集群的 Zadig 环境
- 已安装并配置好的容器镜像仓库（如 Harbor、ACR、ECR 等）
- 已配置 Git 代码仓库集成（如 GitLab、GitHub 等）
- 代码仓库已 Fork 或克隆到你的账户中

## 1. 准备代码仓库

### 1.1 项目结构

整个项目的目录结构如下：

```
k8s_sample/
├── services/          # 所有服务代码
│   ├── backend/       # Go 后端服务 (main.go, go.mod, Dockerfile)
│   ├── frontend/      # UmiJS 前端源代码
│   └── nginx/         # Nginx Dockerfile + default.conf（使用前端构建产物）
├── build/             # 构建脚本（用于 Zadig 配置）
│   ├── backend.sh     # 后端构建脚本
│   └── nginx.sh       # 前端构建脚本
└── deploy/
    └── k8s/           # Kubernetes YAML 文件（分离结构）
        ├── namespace.yaml
        ├── backend/
        │   ├── deployment.yaml
        │   └── service.yaml
        └── nginx/
            ├── deployment.yaml
            ├── service.yaml
            └── ingress.yaml
```

### 1.2 配置代码源

在 Zadig 中配置 Git 代码仓库集成：

1. 进入 Zadig 系统设置
2. 选择"代码源" → "新增代码源"
3. 选择你的 Git 平台（GitLab/GitHub 等）
4. 配置认证信息（Application ID 和 Secret）
5. 完成授权

> 详细步骤请参考 [Zadig 代码源集成文档](https://docs.koderover.com/)

## 2. 创建 K8s YAML 项目

### 2.1 新建项目

1. 在 Zadig 项目模块中，点击"新建项目"
2. 填写项目基本信息：
   - **项目名称**：`micro-demo`（或自定义）
   - **项目描述**：微服务演示项目
   - **项目类型**：**选择 "K8s YAML 项目"** ⚠️ 重要
3. 点击"立即新建"按钮

### 2.2 从代码库同步服务定义

Zadig 会自动识别项目中的服务定义。操作步骤：

1. 在项目初始向导中，点击"从代码库同步"按钮
2. 配置代码源信息：
   - **代码源**：选择你配置的代码源
   - **组织名/用户名**：选择你的账户名
   - **代码库**：选择 `k8s_sample`（或你的仓库名）
   - **分支**：选择 `main`（或你的主分支）
3. **选择文件(夹)**：点击"选择文件夹"按钮
   - 第一次选择：`deploy/k8s/backend`（识别后端服务）
   - 第二次选择：`deploy/k8s/nginx`（识别前端服务）
4. 点击"同步"按钮，Zadig 将自动识别并导入服务定义

> **说明**：Zadig 会自动识别目录下的 Deployment、Service、Ingress 等资源。每个服务目录对应一个服务组件。

### 2.3 验证服务识别

同步完成后，你应该看到两个服务：

- **go-backend**：后端服务（包含 deployment.yaml 和 service.yaml）
- **web-entry**：前端入口服务（包含 deployment.yaml、service.yaml 和 ingress.yaml）

点击服务名称可以查看相关的定义文件。

> **关于镜像地址**：
> - YAML 文件中的镜像地址（如 `registry.example.com/micro-demo/go-backend:v1.0.0`）是示例值
> - 在 Zadig 中，构建脚本使用 `$IMAGE` 变量会自动更新镜像地址
> - 如果手动部署，需要将示例地址替换为实际的镜像仓库地址

## 3. 配置服务构建

### 3.1 后端服务构建配置

点击 **go-backend** 服务名称，添加构建配置：

1. **依赖的软件包**：选择系统内置的 Go 版本（如 go 1.21）
2. **代码信息**：
   - 选择你的代码库
   - 默认分支选择 `main`
3. **构建脚本**：可以直接复制 `build/backend.sh` 文件中的内容，或填入以下内容：

```bash
#!/bin/bash
set -e

cd services/backend
make build                    # 使用 Makefile 构建（推荐）

docker build -t $IMAGE -f Dockerfile .
docker push $IMAGE
```

> **提示**：构建脚本已保存在 `build/backend.sh`，可以直接复制使用。

> **说明**：也可以直接使用命令 `go mod tidy && go build -o server .`，但使用 Makefile 更符合官方示例的标准。

> **重要**：必须使用 `$IMAGE` 变量，这是 Zadig 自动提供的环境变量，包含完整的镜像地址（仓库+标签）。

### 3.2 前端服务构建配置

点击 **web-entry** 服务名称，添加构建配置：

1. **代码信息**：
   - 选择你的代码库
   - 默认分支选择 `main`
2. **构建脚本**：可以直接复制 `build/nginx.sh` 文件中的内容，或填入以下内容：

```bash
#!/bin/bash
set -e

docker build -t $IMAGE -f services/nginx/Dockerfile .
docker push $IMAGE
```

> **提示**：构建脚本已保存在 `build/nginx.sh`，可以直接复制使用。

> **说明**：`services/nginx/Dockerfile` 会在构建时自动处理前端代码的编译和打包，无需单独的前端构建步骤。

### 3.3 容器端口配置

- **go-backend**：容器端口 `8080`
- **web-entry**：容器端口 `80`

## 4. 配置环境

### 4.1 创建环境

在项目配置向导中，配置环境：

1. **环境配置**：
   - 创建两个测试环境：`dev` 和 `qa`
   - 生产环境 `prod` 可以稍后创建
2. **访问入口配置**：
   - 为每个环境配置访问入口域名后缀
   - 例如：`edu.koderover.com` 或 `micro-demo.local`
   - 这将用于自动生成 Ingress 访问地址
3. 点击"创建环境"按钮

### 4.2 环境说明

Zadig 会为项目创建默认环境和工作流：

- **dev 环境**：开发环境，用于日常开发测试
- **qa 环境**：测试环境，用于集成测试
- **prod 环境**：生产环境（需要手动创建）

每个环境都有对应的默认工作流：
- **dev 工作流**：部署到开发环境
- **qa 工作流**：部署到测试环境
- **ops 工作流**：部署到生产环境

## 5. 工作流配置

### 5.1 默认工作流

Zadig 会自动创建默认工作流，包含：

1. **构建阶段**：自动构建服务镜像
2. **部署阶段**：自动部署到对应环境

### 5.2 添加测试任务（可选）

可以在工作流中添加测试任务：

1. 进入"测试"标签页
2. 点击"新建测试"按钮
3. 配置测试套件：
   - **测试名称**：如 `unit-test`、`int-test`、`perf-test`
   - **代码信息**：选择项目代码库
   - **测试脚本**：编写测试脚本
   - **测试报告目录**：指定测试结果输出路径
4. 在工作流编辑中添加测试阶段：
   - 找到对应的工作流（如 qa 工作流）
   - 点击编辑按钮
   - 添加测试阶段，并添加测试任务
   - 选择需要执行的测试套件

### 5.3 配置部署任务

对于生产环境（ops 工作流），需要配置部署任务：

1. 找到 **ops 工作流**，点击"编辑"按钮
2. 配置部署任务：
   - **环境**：选择 `prod` 环境
   - **部署内容**：勾选"服务镜像"、"服务变量"和"服务配置"
3. 保存工作流

## 6. 执行工作流

### 6.1 手动触发工作流

1. 进入"工作流"标签页
2. 找到对应的工作流（如 dev 工作流）
3. 点击"执行"按钮
4. 选择要构建和部署的服务（backend 和 web-entry）
5. 点击"执行"开始构建和部署

### 6.2 自动触发工作流

配置代码合并自动触发：

1. 在工作流设置中，配置 Webhook 触发
2. 当代码合并到主分支时，自动触发工作流执行

### 6.3 查看执行结果

工作流执行完成后：

1. 查看构建日志，确认镜像构建成功
2. 查看部署日志，确认服务部署成功
3. 在"环境"标签页查看服务运行状态

## 7. 访问应用

### 7.1 查看访问地址

1. 进入"环境"标签页
2. 选择对应的环境（如 dev）
3. 找到 **web-entry** 服务
4. 点击服务入口，查看访问地址

访问地址格式：`http://<service-name>.<namespace>.<domain>`

例如：`http://web-entry.micro-demo.edu.koderover.com`

### 7.2 本地访问配置

如果使用本地域名（如 `micro.demo.local`）：

1. 查看 Ingress Controller 的 IP 地址
2. 在本地 hosts 文件中添加：
   ```
   <ingress-ip> micro.demo.local
   ```
3. 访问 `http://micro.demo.local/`

### 7.3 登录应用

- 默认账号：`admin`
- 默认密码：`admin`

## 8. 版本管理

### 8.1 创建版本

基于测试环境创建正式版本：

1. 进入"版本管理"标签页
2. 点击"创建版本"按钮
3. 填写版本信息：
   - **版本号**：如 `v1.0.0`
   - **标签**：版本标签
   - **版本描述**：版本说明
4. 选择版本来源：
   - **环境**：选择 `qa` 环境
   - **服务**：选择 `frontend` 和 `backend` 服务
5. 配置镜像版本：
   - 选择生产环境使用的镜像仓库
   - 配置服务的镜像版本（如 `v1`）
6. 点击"完成"创建版本

版本创建后，会包含：
- 服务的镜像信息
- 服务的 YAML 配置信息
- 用于后续追溯和回滚

### 8.2 发布版本到生产环境

1. 找到 **ops 工作流**，点击"执行"按钮
2. 点击"选择版本"
3. 选择先前创建的版本（如 `v1.0.0`）
4. 关闭「服务过滤」，点击"确定"
5. 配置访问入口地址（如 `edu.koderover.com`）
6. 点击"执行"按钮，开始生产发布

### 8.3 版本回滚

如果需要回滚到之前的版本：

1. 进入"环境"标签页
2. 选择生产环境
3. 找到需要回滚的服务（如 backend）
4. 点击服务右侧的"历史版本"
5. 在历史版本列表中选择一个可用版本
6. 点击"回滚"按钮进行回滚操作

## 9. 服务调试

### 9.1 查看服务状态

1. 进入"环境"标签页
2. 选择对应的环境
3. 查看服务的运行状态、运行版本等信息

### 9.2 查看实时日志

1. 点击服务名称（如 go-backend）
2. 点击"实时日志"标签
3. 查看服务的输出日志，分析问题

### 9.3 容器调试

1. 点击服务名称
2. 点击"调试"按钮
3. 进入容器内部，可以执行命令分析网络、磁盘等问题

示例调试命令：

```bash
# 安装工具
apt install curl -y

# 测试后端 API
curl localhost:8080/healthz

# 测试前端
curl localhost:80/
```

## 10. 常见问题

### 10.1 镜像拉取失败

- 检查镜像仓库地址是否正确
- 检查镜像仓库访问权限
- 如果使用私有仓库，需要在 Zadig 中配置 `imagePullSecrets`

### 10.2 服务无法访问

- 检查 Service 端点：`kubectl get endpoints -n micro-demo`
- 检查 Pod 状态：`kubectl get pods -n micro-demo`
- 检查 Ingress 配置：`kubectl get ingress -n micro-demo`

### 10.3 Ingress 无法访问

- 确认集群已安装 Ingress Controller
- 检查 IngressClass：`kubectl get ingressclass`
- 检查 Ingress 状态：`kubectl describe ingress web-entry -n micro-demo`

### 10.4 构建失败

- 检查构建脚本中的路径是否正确
- 确认使用了 `$IMAGE` 变量
- 查看构建日志，定位具体错误

## 11. 高级功能

### 11.1 多环境推广

可以配置多条工作流，实现多环境推广：

- **Dev → QA → Prod**：从开发环境逐步推广到生产环境
- 通过镜像复用或 Helm/Kustomize 参数差异化部署

### 11.2 灰度发布

在生产环境中，可以配置灰度发布策略：

1. 配置金丝雀部署
2. 逐步增加流量比例
3. 验证通过后全量发布

### 11.3 监控和告警

- 接入 Prometheus 进行监控
- 配置告警规则
- 查看服务指标和日志

## 12. 参考资源

- [Zadig 官方文档](https://docs.koderover.com/)
- [Zadig 官方教程](https://koderover.com/tutorials-detail/codelabs/t100/index.html)
- [Zadig 社区论坛](https://community.koderover.com/)

## 总结

通过以上步骤，你已经完成了：

1. ✅ 创建 K8s YAML 项目
2. ✅ 从代码库同步服务定义
3. ✅ 配置服务构建脚本（使用 `$IMAGE` 变量）
4. ✅ 配置环境和访问入口
5. ✅ 配置工作流和测试任务
6. ✅ 执行工作流部署服务
7. ✅ 创建版本并发布到生产环境
8. ✅ 使用调试功能排查问题

整个项目现在可以通过 Zadig 完成自动构建、测试、发布与迭代。祝部署顺利！🎉
