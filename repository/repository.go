package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ranp9275-sketch/blog-backend-golang/models"
	"github.com/redis/go-redis/v9"
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

// ==================== User ====================

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
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
	// 需要通过关联表查询
	var articles []models.Article
	var total int64

	// Count
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

func (r *Repository) GetCommentsByArticleID(articleID string) ([]models.Comment, error) {
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

// ==================== Stats ====================

func (r *Repository) RecordView(view *models.ArticleView) error {
	return r.db.Create(view).Error
}

func (r *Repository) IncrementArticleViews(articleID string) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (r *Repository) GetArticleStats(articleID string) (int64, int64, error) {
	var views int64
	var comments int64

	// 获取浏览量
	var article models.Article
	if err := r.db.Select("views").Where("id = ?", articleID).First(&article).Error; err != nil {
		return 0, 0, err
	}
	views = article.Views

	// 获取评论数
	r.db.Model(&models.Comment{}).Where("article_id = ? AND status = ?", articleID, "approved").Count(&comments)

	return views, comments, nil
}

// ==================== Search ====================

func (r *Repository) SearchArticles(query string, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	searchQuery := "%" + query + "%"

	r.db.Model(&models.Article{}).
		Where("status = ? AND (title LIKE ? OR content LIKE ? OR excerpt LIKE ?)", "published", searchQuery, searchQuery, searchQuery).
		Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("status = ? AND (title LIKE ? OR content LIKE ? OR excerpt LIKE ?)", "published", searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(pageSize).
		Order("published_at DESC").
		Find(&articles).Error

	return articles, total, err
}

// ==================== User ====================

func (r *Repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

// ==================== Favorite ====================

func (r *Repository) GetFavorites(userID string) ([]models.Article, error) {
	var favorites []models.Favorite
	err := r.db.Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	if len(favorites) == 0 {
		return []models.Article{}, nil
	}

	var articleIDs []string
	for _, f := range favorites {
		articleIDs = append(articleIDs, f.ArticleID)
	}

	var articles []models.Article
	err = r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("id IN ?", articleIDs).
		Find(&articles).Error

	return articles, err
}

func (r *Repository) AddFavorite(favorite *models.Favorite) error {
	// 检查是否已收藏
	var existing models.Favorite
	err := r.db.Where("user_id = ? AND article_id = ?", favorite.UserID, favorite.ArticleID).First(&existing).Error
	if err == nil {
		return nil // 已存在，不重复添加
	}
	return r.db.Create(favorite).Error
}

func (r *Repository) RemoveFavorite(userID, articleID string) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{}).Error
}

func (r *Repository) IsFavorited(userID, articleID string) bool {
	var favorite models.Favorite
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite).Error
	return err == nil
}

// ==================== Article Extended ====================

func (r *Repository) GetArticleByIDWithoutStatus(id string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("id = ?", id).
		First(&article).Error
	return &article, err
}

func (r *Repository) GetArticlesByAuthor(authorID string, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	r.db.Model(&models.Article{}).Where("author_id = ?", authorID).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Preload("Author").Preload("Tags").
		Where("author_id = ?", authorID).
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&articles).Error

	return articles, total, err
}

func (r *Repository) GetAllArticles(page, pageSize int, status string) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	dataQuery := r.db.Preload("Category").Preload("Author").Preload("Tags")
	if status != "" {
		dataQuery = dataQuery.Where("status = ?", status)
	}
	err := dataQuery.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&articles).Error

	return articles, total, err
}

func (r *Repository) CreateArticleWithTags(article *models.Article, tagIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		if len(tagIDs) > 0 {
			var tags []models.Tag
			if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) UpdateArticleWithTags(id string, updates map[string]interface{}, tagIDs []string) error {
	// 清除缓存
	cacheKey := fmt.Sprintf("article:%s", id)
	r.redis.Del(context.Background(), cacheKey)

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if tagIDs != nil {
			var article models.Article
			if err := tx.Where("id = ?", id).First(&article).Error; err != nil {
				return err
			}

			var tags []models.Tag
			if len(tagIDs) > 0 {
				if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(&article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})
}

// ==================== User Management ====================

func (r *Repository) UpdateUser(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) GetAllUsers(page, pageSize int, query string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.db.Model(&models.User{})

	if query != "" {
		searchQuery := "%" + query + "%"
		db = db.Where("name LIKE ? OR email LIKE ?", searchQuery, searchQuery)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

func (r *Repository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
