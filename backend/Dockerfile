# backend/Dockerfile
FROM golang:1.23.6-alpine AS dev

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/app

# 安装必要工具
RUN apk add --no-cache git curl bash

COPY go.mod go.sum ./
RUN go mod download

# 安装 Air 以支持热重载
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# 复制源码
COPY . .

# 重载
CMD ["air"]