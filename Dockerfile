# 使用多阶段构建
# 第一阶段：构建阶段
FROM golang:1.20.4 AS builder

WORKDIR /app

# 先复制 go.mod 和 go.sum 文件并下载依赖，利用 Docker 缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建 client 和 server
RUN go build -o client client.go
RUN go build -o server server.go

# 第二阶段：运行阶段
FROM ubuntu:18.04

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/client ./client
COPY --from=builder /app/server ./server

EXPOSE 8080

# 设置默认命令
CMD ["ls", "-l", "/app"]