package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ranp9275-sketch/blog-backend-golang/models"
	"github.com/ranp9275-sketch/blog-backend-golang/repository"
	"golang.org/x/crypto/bcrypt"
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
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
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
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
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
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handlers) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理可能为空的指针字段
	if article.CategoryID != nil && *article.CategoryID == "" {
		article.CategoryID = nil
	}
	if article.AuthorID != nil && *article.AuthorID == "" {
		article.AuthorID = nil
	}

	// 设置作者 ID
	userID, exists := c.Get("userID")
	if exists {
		authorID := userID.(string)
		article.AuthorID = &authorID
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

	// 处理空字符串转为 NULL
	if val, ok := updates["category_id"]; ok {
		if s, ok := val.(string); ok && s == "" {
			updates["category_id"] = nil
		}
	}
	if val, ok := updates["author_id"]; ok {
		if s, ok := val.(string); ok && s == "" {
			updates["author_id"] = nil
		}
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

	comments, err := h.repo.GetCommentsByArticleID(articleID)
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

	// 增加文章浏览次数
	h.repo.IncrementArticleViews(articleID)

	c.JSON(http.StatusOK, gin.H{"message": "View recorded"})
}

func (h *Handlers) GetArticleStats(c *gin.Context) {
	articleID := c.Param("id")

	views, comments, err := h.repo.GetArticleStats(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"views": views, "comments": comments})
}

// ==================== Stats Handlers ====================

// func (h *Handlers) GetStats(c *gin.Context) {
// 	stats, err := h.repo.GetStats()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, stats)
// }

// ==================== Search Handlers ====================

func (h *Handlers) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	articles, total, err := h.repo.SearchArticles(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// ==================== User Handlers ====================

func (h *Handlers) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, err := h.repo.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ==================== Favorite Handlers ====================

func (h *Handlers) GetFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	favorites, err := h.repo.GetFavorites(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *Handlers) AddFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		ArticleID string `json:"article_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	favorite := &models.Favorite{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		ArticleID: req.ArticleID,
		CreatedAt: time.Now(),
	}

	if err := h.repo.AddFavorite(favorite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Favorite added"})
}

func (h *Handlers) RemoveFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	articleID := c.Param("articleID")

	if err := h.repo.RemoveFavorite(userID.(string), articleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed"})
}

// ==================== User Article Handlers ====================

// CreateUserArticle 用户创建文章
func (h *Handlers) CreateUserArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Title      string      `json:"title" binding:"required"`
		Content    string      `json:"content" binding:"required"`
		Excerpt    string      `json:"excerpt"`
		CoverImage string      `json:"cover_image"`
		CategoryID string      `json:"category_id"`
		TagIDs     []string    `json:"tag_ids"`
		Submit     interface{} `json:"submit"` // 兼容 bool 或其他类型
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "draft"
	var publishedAt *time.Time
	if submit, ok := req.Submit.(bool); ok && submit {
		status = "published"
		now := time.Now()
		publishedAt = &now
	}

	authorID := userID.(string)
	categoryID := req.CategoryID
	var categoryIDPtr *string
	if categoryID != "" {
		categoryIDPtr = &categoryID
	}

	article := &models.Article{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		CategoryID:  categoryIDPtr,
		AuthorID:    &authorID,
		Status:      status,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.CreateArticleWithTags(article, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetUserArticles 获取用户自己的文章
func (h *Handlers) GetUserArticles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	query := c.Query("q")
	status := c.Query("status")
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")

	articles, total, err := h.repo.GetArticlesByAuthor(userID.(string), page, pageSize, query, status, categoryID, tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateUserArticle 用户更新自己的文章
func (h *Handlers) UpdateUserArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	articleID := c.Param("id")

	// 检查文章是否属于该用户
	article, err := h.repo.GetArticleByIDWithoutStatus(articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if article.AuthorID == nil || *article.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own articles"})
		return
	}

	var req struct {
		Title      string      `json:"title"`
		Content    string      `json:"content"`
		Excerpt    string      `json:"excerpt"`
		CoverImage string      `json:"cover_image"`
		CategoryID string      `json:"category_id"`
		TagIDs     []string    `json:"tag_ids"`
		Submit     interface{} `json:"submit"` // 兼容 bool 或其他类型
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Excerpt != "" {
		updates["excerpt"] = req.Excerpt
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.CategoryID != "" {
		updates["category_id"] = req.CategoryID
	}

	// 处理 submit 字段
	if submit, ok := req.Submit.(bool); ok && submit {
		updates["status"] = "published"
		now := time.Now()
		updates["published_at"] = now
	}

	if err := h.repo.UpdateArticleWithTags(articleID, updates, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated"})
}

// DeleteUserArticle 用户删除自己的文章
func (h *Handlers) DeleteUserArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	articleID := c.Param("id")

	// 检查文章是否属于该用户
	article, err := h.repo.GetArticleByIDWithoutStatus(articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if article.AuthorID == nil || *article.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own articles"})
		return
	}

	if err := h.repo.DeleteArticle(articleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}

// ==================== Admin Article Handlers ====================

// GetAllArticlesAdmin 管理员获取所有文章
func (h *Handlers) GetAllArticlesAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status") // 可选：draft, pending, published, rejected

	articles, total, err := h.repo.GetAllArticles(page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetPendingArticles 获取待审核文章
func (h *Handlers) GetPendingArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	articles, total, err := h.repo.GetAllArticles(page, pageSize, "pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// RejectArticle 拒绝文章
func (h *Handlers) RejectArticle(c *gin.Context) {
	articleID := c.Param("id")

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.repo.UpdateArticle(articleID, map[string]interface{}{
		"status":     "rejected",
		"updated_at": time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article rejected"})
}

// ==================== Profile Handlers ====================

// UpdateProfile 更新个人信息
func (h *Handlers) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
		Bio    string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	updates["bio"] = req.Bio

	if err := h.repo.UpdateUser(userID.(string), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// UpdatePassword 修改密码
func (h *Handlers) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户
	user, err := h.repo.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前密码错误"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.repo.UpdateUser(userID.(string), map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

// ==================== Admin User Handlers ====================

// GetAllUsers 获取所有用户
func (h *Handlers) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	query := c.Query("q")

	users, total, err := h.repo.GetAllUsers(page, pageSize, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateUserRole 更新用户角色
func (h *Handlers) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	if err := h.repo.UpdateUser(userID, map[string]interface{}{
		"role":       req.Role,
		"updated_at": time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}

// DeleteUser 删除用户
func (h *Handlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 不能删除自己
	currentUserID, _ := c.Get("userID")
	if userID == currentUserID.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	if err := h.repo.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ==================== Donation QRCode Handlers ====================

// GetDonationQRCodes 获取启用的打赏二维码（公开）
func (h *Handlers) GetDonationQRCodes(c *gin.Context) {
	qrcodes, err := h.repo.GetDonationQRCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qrcodes)
}

// GetAllDonationQRCodes 获取所有打赏二维码（管理员）
func (h *Handlers) GetAllDonationQRCodes(c *gin.Context) {
	qrcodes, err := h.repo.GetAllDonationQRCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qrcodes)
}

// CreateDonationQRCode 创建打赏二维码
func (h *Handlers) CreateDonationQRCode(c *gin.Context) {
	var qrcode models.DonationQRCode
	if err := c.ShouldBindJSON(&qrcode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qrcode.ID = uuid.New().String()
	qrcode.CreatedAt = time.Now()
	qrcode.UpdatedAt = time.Now()

	if err := h.repo.CreateDonationQRCode(&qrcode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qrcode)
}

// UpdateDonationQRCode 更新打赏二维码
func (h *Handlers) UpdateDonationQRCode(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()

	if err := h.repo.UpdateDonationQRCode(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QRCode updated"})
}

// DeleteDonationQRCode 删除打赏二维码
func (h *Handlers) DeleteDonationQRCode(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteDonationQRCode(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QRCode deleted"})
}

// ==================== File Upload Handlers ====================

// UploadFile 上传文件
func (h *Handlers) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件大小 (最大 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// 创建上传目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
	filepath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回文件URL
	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"url":      url,
		"filename": filename,
	})
}

// UploadArticle 上传文章 (PDF/MD)
func (h *Handlers) UploadArticle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件大小 (最大 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".md" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only PDF and Markdown (.md) are allowed"})
		return
	}

	// 获取作者ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 生成文章标题 (去除后缀)
	title := strings.TrimSuffix(file.Filename, ext)
	content := ""

	if ext == ".md" {
		// 读取 MD 内容
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()

		buf := new(strings.Builder)
		if _, err := io.Copy(buf, f); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			return
		}
		content = buf.String()
	} else if ext == ".pdf" {
		// 保存 PDF 文件
		uploadDir := "./uploads/articles"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
		dst := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDF file"})
			return
		}

		// 生成 PDF 查看链接
		pdfURL := fmt.Sprintf("/uploads/articles/%s", filename)

		// 尝试提取 PDF 文本内容
		extractedText, err := extractPDFText(dst)
		if err == nil && len(extractedText) > 0 {
			// 成功提取文本，创建包含文本内容的文章
			content = fmt.Sprintf("<!-- PDF_URL: %s -->\n\n## 📄 查看完整PDF\n\n本文档已上传为PDF格式，您可以：\n- 快速浏览下方提取的文本内容\n- [点击此处在线预览PDF](%s)（保留完整格式和图片）\n- [下载PDF文件](%s)\n\n---\n\n## 文档内容\n\n%s\n\n---\n\n**📎 附件**: [下载原始PDF文件](%s)", pdfURL, pdfURL, pdfURL, extractedText, pdfURL)
		} else {
			// 提取失败，使用下载链接
			content = fmt.Sprintf("<!-- PDF_URL: %s -->\n\n### 文章已作为 PDF 上传\n\n[点击此处在线预览PDF](%s)\n\n[下载PDF文件](%s)\n\n*注：PDF文本提取失败，请使用上方链接查看完整内容*", pdfURL, pdfURL, pdfURL)
		}
	}

	// 创建文章对象
	now := time.Now()
	authorID := userID.(string)

	article := &models.Article{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		Status:      "published",
		AuthorID:    &authorID,
		CategoryID:  nil, // 不设置默认分类，允许为空
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: &now,
	}

	// 保存到数据库
	if err := h.repo.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}
