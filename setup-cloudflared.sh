#!/bin/bash

# 让 .env 中的变量生效
set -a
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/.env"
set +a

# 确保 cloudflared 已安装
if ! command -v cloudflared &> /dev/null; then
    echo "cloudflared 未安装，请先安装 Cloudflare Tunnel CLI."
    exit 1
fi

# 确保 TUNNEL_NAME 存在
if [[ -z "$TUNNEL_NAME" ]]; then
    echo "TUNNEL_NAME 未在 .env 中定义，请检查 .env 文件."
    exit 1
fi

# 检查 Cloudflare 是否已经登录
if [[ ! -f ~/.cloudflared/cert.pem ]]; then
    echo "未登录 Cloudflare，正在登录..."
    cloudflared tunnel login
else
    echo "已登录 Cloudflare，跳过登录步骤."
fi

# 更新 .env 变量的函数（使用 awk）
update_env_variable() {
    local key="$1"
    local value="$2"
    local env_file="$SCRIPT_DIR/.env"

    awk -v key="$key" -v value="$value" '
    BEGIN { updated = 0 }
    {
        if ($0 ~ "^" key "=") {
            print key "=" value
            updated = 1
        } else {
            print $0
        }
    }
    END {
        if (updated == 0) {
            print key "=" value
        }
    }' "$env_file" > "$env_file.tmp" && mv "$env_file.tmp" "$env_file"
}

# 获取现有 TUNNEL_ID
if [[ -z "$TUNNEL_ID" ]]; then
    echo "检查是否已有同名隧道..."
    TUNNEL_ID=$(cloudflared tunnel list | grep "$TUNNEL_NAME" | awk '{print $1}')
    
    if [[ -n "$TUNNEL_ID" ]]; then
        echo "找到已存在的隧道: $TUNNEL_NAME (ID: $TUNNEL_ID)"
        update_env_variable "TUNNEL_ID" "$TUNNEL_ID"
    else
        echo "未找到已有隧道，创建新的 Cloudflare Tunnel..."
        TUNNEL_ID=$(cloudflared tunnel create "$TUNNEL_NAME" | head -n 1 | grep -oE "^[a-f0-9-]{36}$")

        if [[ -z "$TUNNEL_ID" ]]; then
            echo "创建隧道失败，请检查 Cloudflare 账户权限."
            exit 1
        fi

        echo "隧道创建成功，ID: $TUNNEL_ID"
        update_env_variable "TUNNEL_ID" "$TUNNEL_ID"
    fi
else
    echo "TUNNEL_ID 已在 .env 文件中定义: $TUNNEL_ID"
fi

# 配置 DNS 路由
configure_dns() {
    local domain=$1
    local existing_route=$(cloudflared tunnel route ip list | grep "$domain" | awk '{print $1}')
    
    if [[ -n "$existing_route" ]]; then
        echo "域名 $domain 已存在，跳过 DNS 配置."
    else
        echo "配置隧道的 DNS: $domain"
        cloudflared tunnel route dns "$TUNNEL_NAME" "$domain"
    fi
}

if [[ -n "$TUNNEL_DOMAIN" ]]; then
    configure_dns "$TUNNEL_DOMAIN"
fi

if [[ -n "$TUNNEL_API_DOMAIN" ]]; then
    configure_dns "$TUNNEL_API_DOMAIN"
fi

# 复制凭据文件
mkdir -p "$SCRIPT_DIR/cloudflared"
CREDENTIALS_FILE="$SCRIPT_DIR/cloudflared/credentials.json"

if [[ -f ~/.cloudflared/"$TUNNEL_ID".json ]]; then
    cp ~/.cloudflared/"$TUNNEL_ID".json "$CREDENTIALS_FILE"
else
    echo "错误: 找不到 ~/.cloudflared/$TUNNEL_ID.json，凭据文件复制失败."
    exit 1
fi

echo "凭据文件已复制."

# 获取 Tunnel Token 并更新 `.env`
echo "获取 Cloudflare Tunnel Token..."
TUNNEL_TOKEN=$(cloudflared tunnel token "$TUNNEL_ID")

if [[ -z "$TUNNEL_TOKEN" ]]; then
    echo "获取 TUNNEL_TOKEN 失败，请检查 Cloudflare 设置."
    exit 1
fi

update_env_variable "TUNNEL_TOKEN" "$TUNNEL_TOKEN"
echo ".env 文件已更新 TUNNEL_TOKEN"

# 生成最终的 config.yml
envsubst < "$SCRIPT_DIR/cloudflared/config.template.yml" > "$SCRIPT_DIR/cloudflared/config.yml"
echo "配置文件已更新: config.yml"

echo "Cloudflare Tunnel 配置完成!"
