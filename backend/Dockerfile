# backend/Dockerfile

# 开发环境 Dockerfile (支持热重载)
FROM golang:1.23.6-bullseye AS dev

# 设置国内镜像和模块代理
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0

# 设置工作目录为模块路径
WORKDIR /app

# 安装基础工具链
RUN apt-get update && apt-get install -y \
    git \
    curl \
    bash \
    gcc \
    musl-dev \
    && rm -rf /var/lib/apt/lists/*

# 安装 Air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# 预下载依赖 (利用 Docker 缓存层)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制 Air 配置文件
COPY .air.toml ./

# 复制项目文件 (排除不需要的目录)
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg


# 启动热重载
CMD ["air", "-c", ".air.toml"]