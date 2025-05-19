package models

import "gorm.io/gorm"

var DB *gorm.DB

type ForumCategory struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;not null"`
}

// 预定义的论坛分类
var DefaultCategories = []ForumCategory{
	{ID: 1, Name: "技术讨论"},
	{ID: 2, Name: "学习资源"},
	{ID: 3, Name: "经验分享"},
	{ID: 4, Name: "求助问答"},
	{ID: 5, Name: "活动公告"},
}

// 初始化论坛分类数据
func InitForumCategories() error {
	for _, category := range DefaultCategories {
		if err := DB.FirstOrCreate(&category, "id = ?", category.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
