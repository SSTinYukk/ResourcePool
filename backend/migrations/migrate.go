package migrations

import (
	"log"

	"gorm.io/gorm"

	"g/front/backend/models"
)

// RunMigrations 运行数据库迁移
func RunMigrations(db *gorm.DB) {
	log.Println("开始数据库迁移...")

	// 自动迁移数据库表结构
	err := db.AutoMigrate(
		&models.User{},
		&models.Resource{},
		&models.Category{},
		&models.PointRecord{},
		&models.Topic{},
		&models.Reply{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.Comment{},
		&models.UserFavorite{},
		&models.UserLike{},
		&models.UserTopicFavorite{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")
}
