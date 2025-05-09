package models

import (
	"time"

	"gorm.io/gorm"
)

// UserLike represents a user's like on a resource
type UserLike struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	ResourceID uint           `json:"resource_id"`
	Resource   Resource       `json:"resource" gorm:"foreignKey:ResourceID"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
