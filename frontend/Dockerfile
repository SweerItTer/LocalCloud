# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install --frozen-lockfile
COPY . .
RUN npm run build

# 生产阶段
FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY docker/nginx/nginx.conf /etc/nginx/conf.d/default.conf

# 开发阶段（调试用）
# 使用：docker-compose up frontend
# 访问：localhost:3000 带热更新
FROM node:18-alpine AS dev
WORKDIR /app
CMD ["npm", "start"]