.PHONY: build build-frontend clean help

# 默认目标
.DEFAULT_GOAL := help

# 构建步骤（用于构建脚本，实际构建在 Dockerfile 中完成）
build:
	@echo "Preparing nginx build..."
	@echo "Note: Frontend build will be handled by Dockerfile during docker build"

# 构建前端（在 Dockerfile 构建阶段会自动执行，此 Makefile 用于本地测试）
build-frontend:
	@echo "Building frontend..."
	cd ../frontend && npm install && npm run build

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ../frontend/dist
	rm -rf ../frontend/node_modules

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  make build-frontend - Build frontend application (for local testing)"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help message"
	@echo ""
	@echo "Note: In production, frontend build is handled by Dockerfile."

