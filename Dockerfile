# 第一阶段：构建 Go 二进制
FROM golang:1.24-alpine AS builder

# 安装 git 等构建依赖（如果需要拉依赖）
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 拷贝源码
COPY . .

RUN go mod download
# 构建 Go 应用（Linux amd64，输出为 redirector）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /redirector main.go

# 第二阶段：构建精简运行镜像
FROM alpine:3.20

WORKDIR /app

# 拷贝构建产物
COPY --from=builder /redirector .

# 设置权限（一般非必需，go build 默认可执行）
RUN chmod +x redirector

# 启动入口
CMD ["./redirector"]