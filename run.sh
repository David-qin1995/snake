#!/bin/bash

# 确保必要的目录存在
mkdir -p web/static

# 构建前端
echo "Building frontend..."
cd snake-game-frontend && npm run build && cd ..

# 复制前端文件到静态目录
echo "Copying frontend files to static directory..."
cp -r snake-game-frontend/dist/* web/static/

# 启动服务器
echo "Starting server..."
go run cmd/server/main.go 