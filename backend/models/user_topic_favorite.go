package models

import (
	"time"

	"gorm.io/gorm"
)

type UserTopicFavorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	TopicID   uint           `json:"topic_id"`
	Topic     Topic          `json:"topic" gorm:"foreignKey:TopicID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
