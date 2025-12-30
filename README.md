# 技术博客系统 - Golang Gin 后端

## 📚 项目概述

这是一个基于 **Golang Gin 框架**的高性能技术博客管理系统后端。

### 核心特性

- ✅ **高性能**：使用 Gin 框架，支持高并发
- ✅ **缓存优化**：集成 Redis 缓存，提升响应速度
- ✅ **完整的 API**：30+ RESTful API 接口
- ✅ **数据库**：MySQL 数据库，支持事务
- ✅ **认证系统**：JWT Token 认证
- ✅ **权限控制**：基于角色的访问控制
- ✅ **浏览量统计**：记录文章访问数据
- ✅ **评论系统**：支持文章评论

## 🏗️ 项目结构

```
blog-backend-golang/
├── main.go                 # 主入口
├── go.mod                  # 依赖管理
├── config/
│   └── config.go          # 配置管理
├── models/
│   └── models.go          # 数据模型
├── repository/
│   └── repository.go      # 数据访问层
├── handlers/
│   └── handlers.go        # 请求处理器
├── middleware/
│   └── middleware.go      # 中间件
├── .env.example           # 环境变量示例
├── Dockerfile             # Docker 配置
└── README.md              # 项目说明
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

### 安装依赖

```bash
go mod download
go mod tidy
```

### 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和 Redis 连接
```

### 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## 📡 API 接口

### 公开接口

#### 文章相关
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `GET /api/articles/category/:categoryID` - 按分类获取文章
- `GET /api/articles/tag/:tagID` - 按标签获取文章

#### 分类相关
- `GET /api/categories` - 获取分类列表

#### 标签相关
- `GET /api/tags` - 获取标签列表

#### 评论相关
- `GET /api/articles/:id/comments` - 获取评论
- `POST /api/articles/:id/comments` - 创建评论

#### 浏览量统计
- `POST /api/articles/:id/view` - 记录浏览
- `GET /api/articles/:id/stats` - 获取统计数据

### 受保护接口（需要 JWT Token）

#### 文章管理
- `POST /api/admin/articles` - 创建文章
- `PUT /api/admin/articles/:id` - 更新文章
- `DELETE /api/admin/articles/:id` - 删除文章
- `PATCH /api/admin/articles/:id/publish` - 发布文章

#### 分类管理
- `POST /api/admin/categories` - 创建分类
- `PUT /api/admin/categories/:id` - 更新分类
- `DELETE /api/admin/categories/:id` - 删除分类

#### 标签管理
- `POST /api/admin/tags` - 创建标签
- `PUT /api/admin/tags/:id` - 更新标签
- `DELETE /api/admin/tags/:id` - 删除标签

#### 评论管理
- `DELETE /api/admin/comments/:id` - 删除评论

#### 统计数据
- `GET /api/admin/stats` - 获取系统统计
- `GET /api/admin/stats/articles` - 获取文章统计

## 🔐 认证

使用 JWT Token 进行认证。在请求头中添加：

```
Authorization: Bearer <token>
```

## 🗄️ 数据库模型

### User（用户）
- id: UUID
- email: 邮箱
- name: 名称
- role: 角色（admin/user）
- created_at: 创建时间
- updated_at: 更新时间

### Category（分类）
- id: UUID
- name: 分类名称
- slug: URL 友好名称
- desc: 描述
- created_at: 创建时间
- updated_at: 更新时间

### Tag（标签）
- id: UUID
- name: 标签名称
- slug: URL 友好名称
- created_at: 创建时间
- updated_at: 更新时间

### Article（文章）
- id: UUID
- title: 标题
- slug: URL 友好名称
- content: 内容
- excerpt: 摘要
- cover_image: 封面图片
- category_id: 分类 ID
- status: 状态（draft/published）
- views: 浏览次数
- author_id: 作者 ID
- created_at: 创建时间
- updated_at: 更新时间
- published_at: 发布时间

### Comment（评论）
- id: UUID
- article_id: 文章 ID
- author: 作者名称
- email: 邮箱
- content: 评论内容
- status: 状态（pending/approved/rejected）
- created_at: 创建时间
- updated_at: 更新时间

### ArticleView（浏览记录）
- id: UUID
- article_id: 文章 ID
- ip: IP 地址
- user_agent: 用户代理
- created_at: 创建时间

## 📦 Docker 部署

### 构建镜像

```bash
docker build -t blog-backend-golang .
```

### 运行容器

```bash
docker run -d \
  --name blog-backend \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_USER=root \
  -e DB_PASSWORD=password \
  -e DB_NAME=blog \
  -e REDIS_ADDR=redis:6379 \
  -e JWT_SECRET=your-secret-key \
  blog-backend-golang
```

### Docker Compose

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: blog
    ports:
      - "3306:3306"

  redis:
    image: redis:7.0
    ports:
      - "6379:6379"

  blog-backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: mysql
      DB_USER: root
      DB_PASSWORD: password
      DB_NAME: blog
      REDIS_ADDR: redis:6379
      JWT_SECRET: your-secret-key
    depends_on:
      - mysql
      - redis
```

## 🧪 测试

### 创建文章

```bash
curl -X POST http://localhost:8080/api/admin/articles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Article",
    "content": "Article content...",
    "category_id": "category-id",
    "status": "draft"
  }'
```

### 获取文章列表

```bash
curl http://localhost:8080/api/articles?page=1&pageSize=10
```

### 获取文章详情

```bash
curl http://localhost:8080/api/articles/article-id
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|------|------|--------|
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户 | root |
| DB_PASSWORD | 数据库密码 | 空 |
| DB_NAME | 数据库名称 | blog |
| REDIS_ADDR | Redis 地址 | localhost:6379 |
| JWT_SECRET | JWT 密钥 | your-secret-key |
| PORT | 服务端口 | 8080 |

## 📝 代码示例

### 创建文章

```go
article := &models.Article{
    Title:      "My Article",
    Content:    "Article content...",
    CategoryID: "category-id",
    Status:     "draft",
    AuthorID:   "author-id",
}

if err := repo.CreateArticle(article); err != nil {
    log.Fatal(err)
}
```

### 获取文章

```go
article, err := repo.GetArticleByID("article-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Article: %+v\n", article)
```

### 更新文章

```go
updates := map[string]interface{}{
    "title":  "Updated Title",
    "status": "published",
}

if err := repo.UpdateArticle("article-id", updates); err != nil {
    log.Fatal(err)
}
```

## 🚀 性能优化

1. **缓存**：使用 Redis 缓存热点数据
2. **数据库索引**：为常用查询字段添加索引
3. **连接池**：使用 GORM 的连接池
4. **异步处理**：使用 Goroutine 处理耗时操作

## 📚 相关资源

- [Gin 文档](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/)
- [Redis 文档](https://redis.io/)
- [JWT 文档](https://jwt.io/)

## 📄 许可证

MIT License

---

**项目创建时间**：2024 年 12 月
**最后更新**：2024 年 12 月 30 日
**版本**：1.0.0

祝您的技术博客平台蓬勃发展！🚀
