package utils

import (
	"log"
	"time"

	"gorm.io/gorm"

	"g/front/backend/models"
)

// InitData 初始化数据库基础数据
func InitData(db *gorm.DB) {
	// 初始化管理员账号
	initAdminUser(db)

	// 初始化资源分类
	initCategories(db)
}

// 初始化管理员账号
func initAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		// 创建管理员账号
		admin := models.User{
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  "$2a$10$rGJbFJnG5tokiSNTmgL4SuuQJD5L6VEYRIWBx7mJMIzDYwPUOyixy", // 加密后的 "admin123"
			Points:    1000,
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("创建管理员账号失败: %v\n", result.Error)
		} else {
			log.Println("管理员账号创建成功")
		}
	}
}

// 初始化资源分类
func initCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)

	if count == 0 {
		// 主分类
		mainCategories := []models.Category{
			{
				Name:        "教材资源",
				Description: "微机原理与接口技术相关教材和参考书",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实验资料",
				Description: "实验指导、实验报告和相关资料",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "课件讲义",
				Description: "教师课件和讲义资料",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "习题资料",
				Description: "习题集、作业和答案",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "程序代码",
				Description: "汇编语言程序和示例代码",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		for _, category := range mainCategories {
			result := db.Create(&category)
			if result.Error != nil {
				log.Printf("创建分类失败: %v\n", result.Error)
			} else {
				log.Printf("成功创建分类: %s\n", category.Name)
			}
		}

		// 创建论坛分类
		forumCategories := []models.Category{
			{
				Name:        "课程讨论",
				Description: "微机原理课程相关讨论",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "实验交流",
				Description: "实验过程中的问题和解决方案",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "资源求助",
				Description: "寻找特定学习资源",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "作业互助",
				Description: "作业问题讨论和解答",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		for _, category := range forumCategories {
			result := db.Create(&category)
			if result.Error != nil {
				log.Printf("创建论坛分类失败: %v\n", result.Error)
			} else {
				log.Printf("成功创建论坛分类: %s\n", category.Name)
			}
		}

		log.Printf("基础分类创建完成，共创建了%d个主分类和%d个论坛分类\n", len(mainCategories), len(forumCategories))
	}
}
