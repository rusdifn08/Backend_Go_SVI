package domain

import (
	"time"
)

type ArticleStatus string

const (
	StatusPublish ArticleStatus = "publish"
	StatusDraft   ArticleStatus = "draft"
	StatusThrash  ArticleStatus = "thrash"
)

// Article model mapping to `posts` table
type Article struct {
	ID          int           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string        `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Content     string        `gorm:"column:content;type:text;not null" json:"content"`
	Category    string        `gorm:"column:category;type:varchar(100);not null" json:"category"`
	CreatedDate time.Time     `gorm:"column:created_date;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_date"`
	UpdatedDate time.Time     `gorm:"column:updated_date;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_date"`
	Status      ArticleStatus `gorm:"column:status;type:varchar(100);not null" json:"status"`
}

func (Article) TableName() string {
	return "posts"
}
