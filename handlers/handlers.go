package handlers

import (
	"net/http"
	"strconv"
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
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Excerpt    string   `json:"excerpt"`
		CoverImage string   `json:"cover_image"`
		CategoryID string   `json:"category_id"`
		TagIDs     []string `json:"tag_ids"`
		Submit     bool     `json:"submit"` // true=发布, false=保存草稿
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "draft"
	var publishedAt *time.Time
	if req.Submit {
		status = "published"
		now := time.Now()
		publishedAt = &now
	}

	article := &models.Article{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		AuthorID:    userID.(string),
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

	articles, total, err := h.repo.GetArticlesByAuthor(userID.(string), page, pageSize)
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

	if article.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own articles"})
		return
	}

	var req struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Excerpt    string   `json:"excerpt"`
		CoverImage string   `json:"cover_image"`
		CategoryID string   `json:"category_id"`
		TagIDs     []string `json:"tag_ids"`
		Submit     bool     `json:"submit"` // true=提交审核
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
	if req.Submit {
		updates["status"] = "pending"
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

	if article.AuthorID != userID.(string) {
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
