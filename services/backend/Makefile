.PHONY: build clean test run docker-build help

# 默认目标
.DEFAULT_GOAL := help

# 构建 Go 应用
build:
	@echo "Building Go application..."
	go mod tidy
	go build -o server .

# 运行测试
test:
	@echo "Running tests..."
	go test ./...

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -f server

# 运行应用（本地开发）
run: build
	@echo "Running application..."
	./server

# 构建 Docker 镜像（本地测试）
docker-build:
	@echo "Building Docker image..."
	docker build -t go-backend:latest .

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  make build        - Build Go application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make run          - Build and run application"
	@echo "  make docker-build - Build Docker image"
	@echo "  make help         - Show this help message"

