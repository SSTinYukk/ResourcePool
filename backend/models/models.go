package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Email     string         `json:"email" gorm:"size:100;uniqueIndex"`
	Password  string         `json:"-" gorm:"size:100;not null"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Points    int            `json:"points" gorm:"default:0"`
	Role      string         `json:"role" gorm:"size:20;default:'user'"` // user, admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Resource 资源模型
type Resource struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title" gorm:"size:100;not null"`
	Description    string         `json:"description" gorm:"type:text"`
	CategoryID     uint           `json:"category_id"`
	Category       Category       `json:"category" gorm:"foreignKey:CategoryID"`
	FilePath       string         `json:"file_path" gorm:"size:255"`
	FileSize       int64          `json:"file_size"`
	FileType       string         `json:"file_type" gorm:"size:50"`
	DownloadCount  int            `json:"download_count" gorm:"default:0"`
	PointsRequired int            `json:"points_required" gorm:"default:0"`
	Status         string         `json:"status" gorm:"size:20;default:'pending'"` // pending, approved, rejected
	UserID         uint           `json:"user_id"`
	User           User           `json:"user" gorm:"foreignKey:UserID"`
	Likes          []UserLike     `json:"likes" gorm:"foreignKey:ResourceID"`
	Favorites      []UserFavorite `json:"favorites" gorm:"foreignKey:ResourceID"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// Category 资源分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	ParentID    *uint          `json:"parent_id" gorm:"default:null"`
	Parent      *Category      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	// 非数据库字段，仅用于API响应
	TopicCount int `json:"topic_count,omitempty" gorm:"-"`
	PostCount  int `json:"post_count,omitempty" gorm:"-"`
}

// PointRecord 积分记录模型
type PointRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Points      int       `json:"points"`              // 正为增加，负为减少
	Type        string    `json:"type" gorm:"size:20"` // upload, download, post, reply
	ResourceID  *uint     `json:"resource_id" gorm:"default:null"`
	Resource    *Resource `json:"resource,omitempty" gorm:"foreignKey:ResourceID"`
	Description string    `json:"description" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
}

// Topic 论坛主题模型
type Topic struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"size:100;not null"`
	Content      string         `json:"content" gorm:"type:text;not null"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	CategoryID   uint           `json:"category_id"`
	Category     Category       `json:"category" gorm:"foreignKey:CategoryID"`
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	ReplyCount   int            `json:"reply_count" gorm:"default:0"`
	LikeCount    int64          `json:"like_count" gorm:"default:0"`
	DislikeCount int64          `json:"dislike_count" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Reply 论坛回复模型
type Reply struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	TopicID   uint           `json:"topic_id"`
	Topic     Topic          `json:"topic" gorm:"foreignKey:TopicID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserFavorite struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	ResourceID uint           `json:"resource_id"`
	Resource   Resource       `json:"resource" gorm:"foreignKey:ResourceID"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
