# LocalCloud

一个尝试平替华为云相册的，基于内网穿透在本地部署云端文件管理的项目。

## 项目愿景

LocalCloud 旨在打造一款**轻量化、零成本、易部署**的私有云盘，让用户能够安全、便捷地管理自己的照片和文件。

## 核心优势

* **零服务器成本**: 无需购买云服务器，只需一台闲置电脑或NAS设备即可搭建。
* **私有化部署**: 所有数据都存储在本地，彻底掌控自己的隐私。
* **容器化部署**: 采用 Docker Compose 一键部署，告别繁琐的配置过程。
* **移动端优先**: 专注于照片管理，提供极致的移动端备份体验。
* **持续迭代**:  项目将不断完善，增加更多实用功能。

## 功能规划

### 核心功能分层

| 阶段 | 必须实现 | 优化项 | 延伸功能 |
|---|---|---|---|
| V1.0 | 容器化部署<br>图片上传/下载<br>多用户隔离<br>响应式Web界面 | 缩略图生成<br>EXIF信息保留 | - |
| V2.0 | 文件版本控制<br>分享链接生成<br>手机端PWA应用 | 智能相册分类<br>WebDAV协议支持 | - |
| V3.0 | 跨设备同步引擎<br>端到端加密<br>AppStore/Play商店上架 | AI图片去重<br>NAS设备预装包 | - |

### 技术选型

* **基础设施层**:*
    * **容器化**: Docker Compose
    * **穿透方案**: Cloudflare Tunnel (免费) / Tailscale (可选)
    * **存储引擎**: MinIO (S3兼容) + Redis (缓存)
* **服务端层**:*
    * **语言框架**: Golang + Gin
    * **鉴权体系**: JWT + OAuth2.0
    * **任务队列**: RabbitMQ
* **客户端层**:*
    * **Web端**: React + Material-UI (SSR)
    * **移动端**: Flutter
    * **桌面端**: Electron (未来)

## 快速开始

1. **准备环境**: 安装 Docker, Docker Compose, cloudflared
2. **拉取代码**: `git clone https://github.com/SweerItTer/LocalCloud.git`
3. **配置参数**: 
   1. 配置`.env`的 TUNNEL_NAME, TUNNEL_DOMAIN 和 TUNNEL_API_DOMAIN (如果没有域名可以试试搜索 **USKG**)(如果是通过网页创建的隧道,需要到网站添加配置,需要填写的数据参考`./cloudflare/config.template.yml`,如果是控制台通过命令创建的,直接运行`setup-cloudflared.sh`即可)
   2. 配置`./backend/.env`的 Github OAuth ID 和 数据库 相关配置
4. **启动服务**: `docker-compose --env-file .env up -d`
5. **访问应用**: 在浏览器中输入 `http://localhost`，即可访问网页服务(app服务待添加)

## 参与贡献

欢迎任何形式的贡献，包括但不限于：

* 提交 Bug 报告
* 提出 Feature 请求
* 参与代码开发
* 完善文档

## 许可证

本项目采用 MIT 许可证。

## 联系我们

* GitHub: [https://github.com/SweerItTer/LocalCloud](https://github.com/SweerItTer/LocalCloud)
* Email: sweeritter@gmail.com / xxxzhou_xian@163.com

---

## 感谢

感谢以下项目为 LocalCloud 提供了技术支持：

* MinIO
* Golang
* React
* Flutter
* ...
