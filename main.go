package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ranp9275-sketch/blog-backend-golang/config"
	"github.com/ranp9275-sketch/blog-backend-golang/handlers"
	"github.com/ranp9275-sketch/blog-backend-golang/middleware"
	"github.com/ranp9275-sketch/blog-backend-golang/models"
	"github.com/ranp9275-sketch/blog-backend-golang/repository"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移表
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
		&models.ArticleView{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis
	redisClient := config.InitRedis(cfg)
	defer redisClient.Close()

	// 初始化仓储
	repo := repository.NewRepository(db, redisClient)

	// 创建 Gin 路由
	router := gin.Default()

	// 添加中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// 初始化处理器
	h := handlers.NewHandlers(repo)

	// 公开路由
	public := router.Group("/api")
	{
		// 文章相关
		public.GET("/articles", h.GetArticles)
		public.GET("/articles/:id", h.GetArticleByID)
		public.GET("/articles/category/:categoryID", h.GetArticlesByCategory)
		public.GET("/articles/tag/:tagID", h.GetArticlesByTag)

		// 分类相关
		public.GET("/categories", h.GetCategories)

		// 标签相关
		public.GET("/tags", h.GetTags)

		// 评论相关
		public.GET("/articles/:id/comments", h.GetComments)
		public.POST("/articles/:id/comments", h.CreateComment)

		// 浏览量统计
		public.POST("/articles/:id/view", h.RecordView)
		public.GET("/articles/:id/stats", h.GetArticleStats)
	}

	// 受保护的路由（需要认证）
	protected := router.Group("/api/admin")
	protected.Use(middleware.AuthMiddleware())
	{
		// 文章管理
		protected.POST("/articles", h.CreateArticle)
		protected.PUT("/articles/:id", h.UpdateArticle)
		protected.DELETE("/articles/:id", h.DeleteArticle)
		protected.PATCH("/articles/:id/publish", h.PublishArticle)

		// 分类管理
		protected.POST("/categories", h.CreateCategory)
		protected.PUT("/categories/:id", h.UpdateCategory)
		protected.DELETE("/categories/:id", h.DeleteCategory)

		// 标签管理
		protected.POST("/tags", h.CreateTag)
		protected.PUT("/tags/:id", h.UpdateTag)
		protected.DELETE("/tags/:id", h.DeleteTag)

		// 评论管理
		protected.DELETE("/comments/:id", h.DeleteComment)

		// 统计数据
		protected.GET("/stats", h.GetStats)
		protected.GET("/stats/articles", h.GetArticleStats)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
