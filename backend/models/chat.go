package models

import (
	"time"
)

// ChatSession 聊天会话模型
type ChatSession struct {
	ID        uint          `json:"id" gorm:"primaryKey"`
	UserID    uint          `json:"user_id" gorm:"index"`
	Title     string        `json:"title"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Messages  []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:SessionID"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SessionID uint      `json:"session_id" gorm:"index"`
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role" gorm:"type:varchar(20)"` // user 或 assistant
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
