package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ranp9275-sketch/blog-backend-golang/models"
	"github.com/ranp9275-sketch/blog-backend-golang/repository"
)

type Handlers struct {
	repo *repository.Repository
}

func NewHandlers(repo *repository.Repository) *Handlers {
	return &Handlers{repo: repo}
}

// ==================== Article Handlers ====================

func (h *Handlers) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	articles, total, err := h.repo.GetArticles(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      articles,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

func (h *Handlers) GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	article, err := h.repo.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *Handlers) GetArticlesByCategory(c *gin.Context) {
	categoryID := c.Param("categoryID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	articles, total, err := h.repo.GetArticlesByCategory(categoryID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      articles,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

func (h *Handlers) GetArticlesByTag(c *gin.Context) {
	tagID := c.Param("tagID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	articles, total, err := h.repo.GetArticlesByTag(tagID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      articles,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

func (h *Handlers) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article.ID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	if err := h.repo.CreateArticle(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *Handlers) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()

	if err := h.repo.UpdateArticle(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated"})
}

func (h *Handlers) DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}

func (h *Handlers) PublishArticle(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	if err := h.repo.UpdateArticle(id, map[string]interface{}{
		"status":       "published",
		"published_at": now,
		"updated_at":   now,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article published"})
}

// ==================== Category Handlers ====================

func (h *Handlers) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handlers) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	if err := h.repo.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handlers) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()

	if err := h.repo.UpdateCategory(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

func (h *Handlers) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

// ==================== Tag Handlers ====================

func (h *Handlers) GetTags(c *gin.Context) {
	tags, err := h.repo.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *Handlers) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.ID = uuid.New().String()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	if err := h.repo.CreateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func (h *Handlers) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()

	if err := h.repo.UpdateTag(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag updated"})
}

func (h *Handlers) DeleteTag(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
}

// ==================== Comment Handlers ====================

func (h *Handlers) GetComments(c *gin.Context) {
	articleID := c.Param("id")

	comments, err := h.repo.GetComments(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *Handlers) CreateComment(c *gin.Context) {
	articleID := c.Param("id")
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.ID = uuid.New().String()
	comment.ArticleID = articleID
	comment.Status = "pending"
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := h.repo.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *Handlers) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// ==================== View Handlers ====================

func (h *Handlers) RecordView(c *gin.Context) {
	articleID := c.Param("id")
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	view := &models.ArticleView{
		ArticleID: articleID,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	if err := h.repo.RecordView(view); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View recorded"})
}

func (h *Handlers) GetArticleStats(c *gin.Context) {
	articleID := c.Param("id")

	views, err := h.repo.GetArticleViews(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"views": views})
}

// ==================== Stats Handlers ====================

func (h *Handlers) GetStats(c *gin.Context) {
	stats, err := h.repo.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
