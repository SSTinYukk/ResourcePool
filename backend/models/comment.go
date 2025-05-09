package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 资源评论模型
type Comment struct {
	gorm.Model
	ResourceID string    `gorm:"type:varchar(36);not null;index"` // 资源ID
	UserID     uint      `gorm:"not null;index"`                  // 用户ID
	Rating     int       `gorm:"not null"`                        // 评分(1-5)
	Content    string    `gorm:"type:text;not null"`              // 评论内容
	Time       time.Time `gorm:"not null"`                        // 评论时间

	// 关联模型
	User     User     `gorm:"foreignKey:UserID"`
	Resource Resource `gorm:"foreignKey:ResourceID"`
}
