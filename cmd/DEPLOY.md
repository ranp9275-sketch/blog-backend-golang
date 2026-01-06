# 部署指南

## 准备工作

### 1. 创建证书目录并放入 SSL 证书

```bash
mkdir -p nginx/certs
# 将你的证书文件放入 nginx/certs 目录
# - fullchain.pem (证书链)
# - privkey.pem (私钥)
```

### 2. 创建环境变量文件

```bash
cp .env.example .env
# 编辑 .env 文件，修改密码和密钥
```

## 部署

### 启动所有服务

```bash
docker-compose up -d
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f nginx
```

### 停止服务

```bash
docker-compose down
```

### 重新构建并启动

```bash
docker-compose up -d --build
```

## 创建管理员用户

首次部署后，需要手动将用户设置为管理员：

```bash
# 进入 MySQL 容器
docker exec -it blog-mysql mysql -u blog -p

# 在 MySQL 中执行
USE newblog;
UPDATE users SET role = 'admin' WHERE email = 'your-email@example.com';
```

## 目录结构

```
.
├── docker-compose.yml      # Docker Compose 配置
├── .env                    # 环境变量（需创建）
├── nginx/
│   ├── nginx.conf          # Nginx 配置
│   └── certs/              # SSL 证书目录
│       ├── fullchain.pem   # 证书链
│       └── privkey.pem     # 私钥
├── blog-backend-golang/    # 后端代码
└── blog-frontend-vue/      # 前端代码
```

## 端口说明

- 80: HTTP (自动重定向到 HTTPS)
- 443: HTTPS
- 3306: MySQL (可选，生产环境建议关闭)
- 6379: Redis (可选，生产环境建议关闭)
- 8080: 后端 API (内部)
- 3000: 前端 (内部)

## 常见问题

### 1. 数据库连接失败
等待 MySQL 完全启动后再启动后端服务：
```bash
docker-compose up -d mysql redis
# 等待 30 秒
docker-compose up -d
```

### 2. 证书问题
确保证书文件权限正确：
```bash
chmod 644 nginx/certs/fullchain.pem
chmod 600 nginx/certs/privkey.pem
```

### 3. 前端构建失败
清理并重新构建：
```bash
docker-compose down
docker system prune -f
docker-compose up -d --build
```
