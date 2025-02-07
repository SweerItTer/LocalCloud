# backend/Dockerfile
FROM golang:1.21-alpine AS dev

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

# 安装热重载工具Air
RUN go install github.com/cosmtrek/air@latest

COPY . .
CMD ["air"]  # 启动时自动检测代码变化