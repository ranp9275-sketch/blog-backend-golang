package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // 密码字段，不返回给前端
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Role      string    `json:"role"` // admin, user
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Article 文章模型
type Article struct {
	ID          string     `gorm:"primaryKey;size:191" json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `gorm:"type:longtext" json:"content"`
	Excerpt     string     `json:"excerpt"`
	CoverImage  string     `json:"cover_image"`
	CategoryID  string     `gorm:"size:191" json:"category_id"`
	Category    Category   `gorm:"foreignKey:CategoryID" json:"category"`
	Status      string     `json:"status"` // draft, published
	Views       int64      `json:"views"`
	AuthorID    string     `gorm:"size:191" json:"author_id"`
	Author      User       `gorm:"foreignKey:AuthorID" json:"author"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`

	// 关联
	Tags     []Tag     `gorm:"many2many:article_tags" json:"tags"`
	Comments []Comment `gorm:"foreignKey:ArticleID" json:"comments"`
}

// Comment 评论模型
type Comment struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	ArticleID string    `gorm:"size:191" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	Status    string    `json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleView 文章浏览记录
type ArticleView struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	ArticleID string    `gorm:"size:191" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// Favorite 用户收藏
type Favorite struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	UserID    string    `gorm:"size:191;index" json:"user_id"`
	ArticleID string    `gorm:"size:191;index" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	CreatedAt time.Time `json:"created_at"`
}

// DonationQRCode 打赏二维码
type DonationQRCode struct {
	ID        string    `gorm:"primaryKey;size:191" json:"id"`
	Name      string    `json:"name"`                                          // 支付方式名称：微信、支付宝等
	Icon      string    `json:"icon"`                                          // 图标类型：wechat, alipay, qq
	QRCodeURL string    `gorm:"column:qrcode_url" json:"qrcode_url"`           // 二维码图片URL
	SortOrder int       `gorm:"column:sort_order;default:0" json:"sort_order"` // 排序
	Enabled   bool      `gorm:"default:true" json:"enabled"`                   // 是否启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate 钩子：自动生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (av *ArticleView) BeforeCreate(tx *gorm.DB) error {
	if av.ID == "" {
		av.ID = uuid.New().String()
	}
	return nil
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

func (d *DonationQRCode) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}
