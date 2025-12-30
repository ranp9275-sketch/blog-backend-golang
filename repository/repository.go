package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ranp9275-sketch/blog-backend-golang/models"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
	}
}

// ==================== Article ====================

func (r *Repository) GetArticles(page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	r.db.Model(&models.Article{}).Where("status = ?", "published").Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("status = ?", "published").
		Offset(offset).Limit(pageSize).
		Order("published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

func (r *Repository) GetArticleByID(id string) (*models.Article, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("article:%s", id)
	val, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var article models.Article
		if err := json.Unmarshal([]byte(val), &article); err == nil {
			return &article, nil
		}
	}

	// 从数据库获取
	var article models.Article
	err = r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("id = ? AND status = ?", id, "published").
		First(&article).Error

	if err == nil {
		// 存入缓存
		data, _ := json.Marshal(article)
		r.redis.Set(context.Background(), cacheKey, data, 1*time.Hour)
	}

	return &article, err
}

func (r *Repository) CreateArticle(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *Repository) UpdateArticle(id string, updates map[string]interface{}) error {
	// 清除缓存
	cacheKey := fmt.Sprintf("article:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteArticle(id string) error {
	// 清除缓存
	cacheKey := fmt.Sprintf("article:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Delete(&models.Article{}, "id = ?", id).Error
}

func (r *Repository) GetArticlesByCategory(categoryID string, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	r.db.Model(&models.Article{}).Where("category_id = ? AND status = ?", categoryID, "published").Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("category_id = ? AND status = ?", categoryID, "published").
		Offset(offset).Limit(pageSize).
		Order("published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

func (r *Repository) GetArticlesByTag(tagID string, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	r.db.Model(&models.Article{}).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, "published").
		Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, "published").
		Offset(offset).Limit(pageSize).
		Order("articles.published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// ==================== Category ====================

func (r *Repository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *Repository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *Repository) UpdateCategory(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteCategory(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// ==================== Tag ====================

func (r *Repository) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *Repository) CreateTag(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *Repository) UpdateTag(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteTag(id string) error {
	return r.db.Delete(&models.Tag{}, "id = ?", id).Error
}

// ==================== Comment ====================

func (r *Repository) GetComments(articleID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("article_id = ? AND status = ?", articleID, "approved").
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (r *Repository) CreateComment(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *Repository) DeleteComment(id string) error {
	return r.db.Delete(&models.Comment{}, "id = ?", id).Error
}

// ==================== Views ====================

func (r *Repository) RecordView(view *models.ArticleView) error {
	return r.db.Create(view).Error
}

func (r *Repository) GetArticleViews(articleID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleView{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// ==================== Stats ====================

func (r *Repository) GetStats() (map[string]interface{}, error) {
	var articleCount, categoryCount, tagCount, commentCount int64

	r.db.Model(&models.Article{}).Where("status = ?", "published").Count(&articleCount)
	r.db.Model(&models.Category{}).Count(&categoryCount)
	r.db.Model(&models.Tag{}).Count(&tagCount)
	r.db.Model(&models.Comment{}).Count(&commentCount)

	return map[string]interface{}{
		"articles":   articleCount,
		"categories": categoryCount,
		"tags":       tagCount,
		"comments":   commentCount,
	}, nil
}
