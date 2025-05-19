package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"g/front/backend/models"
)

// ForumController 论坛控制器
type ForumController struct {
	DB       *gorm.DB
	Redis    *redis.Client
	stopChan chan struct{} // 用于停止定时任务的通道
}

// NewForumController 创建论坛控制器实例
func NewForumController(db *gorm.DB, redisClient *redis.Client) *ForumController {
	fc := &ForumController{DB: db, Redis: redisClient, stopChan: make(chan struct{})}
	go fc.syncLikesToDB() // 启动定时同步任务
	return fc
}

// GetCategories 获取论坛分类
func (c *ForumController) GetCategories(ctx *gin.Context) {
	var categories []models.Category
	c.DB.Where("parent_id IS NULL").Find(&categories)

	// 获取每个分类的主题数和帖子数
	for i := range categories {
		// 获取主题数
		var topicCount int64
		c.DB.Model(&models.Topic{}).Where("category_id = ?", categories[i].ID).Count(&topicCount)

		// 获取帖子数（主题+回复）
		var replyCount int64
		c.DB.Model(&models.Reply{}).Where("topic_id IN (SELECT id FROM topics WHERE category_id = ?)", categories[i].ID).Count(&replyCount)

		// 添加到返回结果
		categories[i].TopicCount = int(topicCount)
		categories[i].PostCount = int(topicCount + replyCount)
	}

	ctx.JSON(http.StatusOK, categories)
}

// GetTopics 获取主题列表
func (c *ForumController) GetTopics(ctx *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	categoryID := ctx.Query("category")

	// 构建查询
	query := c.DB.Model(&models.Topic{})

	// 分类过滤
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 执行查询
	var topics []models.Topic
	var total int64

	query.Count(&total)
	query.Preload("User").Preload("Category").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&topics)

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"topics":   topics,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// LikeTopic 点赞主题
func (c *ForumController) DislikeTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// Redis键名
	likeKey := "topic_likes:" + topicID
	dislikeKey := "topic_dislikes:" + topicID
	userLikeKey := "user_likes:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	userDislikeKey := "user_dislikes:" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// 检查用户是否已点赞
	isLiked, err := c.Redis.SIsMember(ctx, userLikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点赞状态失败"})
		return
	}

	// 如果已点赞，先取消点赞
	if isLiked {
		// 调试日志：打印当前Redis键值
		fmt.Printf("[DEBUG] Before remove like - userLikeKey: %s, members: %v\n", userLikeKey, c.Redis.SMembers(ctx, userLikeKey).Val())
		fmt.Printf("[DEBUG] Before remove like - likeKey: %s, value: %v\n", likeKey, c.Redis.Get(ctx, likeKey).Val())

		_, err = c.Redis.SRem(ctx, userLikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
			return
		}
		_, err = c.Redis.Decr(ctx, likeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
			return
		}

		// 调试日志：打印操作后的Redis键值
		fmt.Printf("[DEBUG] After remove like - userLikeKey: %s, members: %v\n", userLikeKey, c.Redis.SMembers(ctx, userLikeKey).Val())
		fmt.Printf("[DEBUG] After remove like - likeKey: %s, value: %v\n", likeKey, c.Redis.Get(ctx, likeKey).Val())
	}

	// 检查用户是否已点踩
	isDisliked, err := c.Redis.SIsMember(ctx, userDislikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点踩状态失败"})
		return
	}

	if isDisliked {
		// 取消点踩
		_, err = c.Redis.SRem(ctx, userDislikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点踩失败"})
			return
		}

		_, err = c.Redis.Decr(ctx, dislikeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点踩数失败"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "取消点踩成功",
			"disliked": false,
			"liked":    false,
			"likes":    c.getLikeCount(topicID),
			"dislikes": c.getDislikeCount(topicID),
		})
	} else {
		// 点踩
		// 调试日志：打印当前Redis键值
		fmt.Printf("[DEBUG] Before add dislike - userDislikeKey: %s, members: %v\n", userDislikeKey, c.Redis.SMembers(ctx, userDislikeKey).Val())
		fmt.Printf("[DEBUG] Before add dislike - dislikeKey: %s, value: %v\n", dislikeKey, c.Redis.Get(ctx, dislikeKey).Val())

		_, err = c.Redis.SAdd(ctx, userDislikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "点踩失败"})
			return
		}

		_, err = c.Redis.Incr(ctx, dislikeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点踩数失败"})
			return
		}

		// 调试日志：打印操作后的Redis键值
		fmt.Printf("[DEBUG] After add dislike - userDislikeKey: %s, members: %v\n", userDislikeKey, c.Redis.SMembers(ctx, userDislikeKey).Val())
		fmt.Printf("[DEBUG] After add dislike - dislikeKey: %s, value: %v\n", dislikeKey, c.Redis.Get(ctx, dislikeKey).Val())

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "点踩成功",
			"disliked": true,
			"liked":    false,
			"likes":    c.getLikeCount(topicID),
			"dislikes": c.getDislikeCount(topicID),
		})
	}

	// 立即同步点踩数据
	c.syncAllDislikes()
}

// AddFavorite 收藏主题
func (c *ForumController) AddFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// 检查主题是否存在
	var topic models.Topic
	if err := c.DB.First(&topic, topicID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 检查是否已收藏
	var existingFavorite models.UserTopicFavorite
	if err := c.DB.Where("user_id = ? AND topic_id = ?", userID, topic.ID).First(&existingFavorite).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "已收藏该主题", "isFavorited": true})
		return
	}

	// 创建收藏记录
	favorite := models.UserTopicFavorite{
		UserID:    userID.(uint),
		TopicID:   topic.ID,
		CreatedAt: time.Now(),
	}

	if err := c.DB.Create(&favorite).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "isFavorited": true, "message": "收藏成功"})
}

// RemoveFavorite 取消收藏主题
func (c *ForumController) RemoveFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// 检查主题是否存在
	var topic models.Topic
	if err := c.DB.First(&topic, topicID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 删除收藏记录
	if err := c.DB.Where("user_id = ? AND topic_id = ?", userID, topic.ID).Delete(&models.UserTopicFavorite{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "已取消收藏", "isFavorited": false})
}

// GetFavoriteStatus 获取主题收藏状态
func (c *ForumController) GetFavoriteStatus(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// 检查收藏状态
	var favorite models.UserTopicFavorite
	if err := c.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&favorite).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"isFavorited": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"isFavorited": true})
}

// syncLikesToDB 定时将Redis中的点赞数据同步到MySQL
func (c *ForumController) syncLikesToDB() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒同步一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.syncAllLikes()
		case <-c.stopChan:
			return
		}
	}
}

// syncAllLikes 同步所有主题的点赞数到数据库
func (c *ForumController) syncAllLikes() {
	// 获取所有主题ID
	var topicIDs []string
	iter := c.Redis.Scan(context.Background(), 0, "topic_likes:*", 0).Iterator()
	for iter.Next(context.Background()) {
		key := iter.Val()
		topicID := strings.TrimPrefix(key, "topic_likes:")
		topicIDs = append(topicIDs, topicID)
	}

	// 同步每个主题的点赞数
	for _, topicID := range topicIDs {
		likeKey := "topic_likes:" + topicID
		dislikeKey := "topic_dislikes:" + topicID

		// 获取点赞数
		likeCount, err := c.Redis.Get(context.Background(), likeKey).Int64()
		if err != nil && err != redis.Nil {
			continue
		}

		// 获取点踩数
		dislikeCount, err := c.Redis.Get(context.Background(), dislikeKey).Int64()
		if err != nil && err != redis.Nil {
			continue
		}

		// 更新数据库
		c.DB.Model(&models.Topic{}).Where("id = ?", topicID).
			Updates(map[string]interface{}{
				"like_count":    likeCount,
				"dislike_count": dislikeCount,
			})
	}
}

func (c *ForumController) syncAllDislikes() {
	// 获取所有主题ID
	var topicIDs []string
	iter := c.Redis.Scan(context.Background(), 0, "topic_dislikes:*", 0).Iterator()
	for iter.Next(context.Background()) {
		key := iter.Val()
		topicID := strings.TrimPrefix(key, "topic_dislikes:")
		topicIDs = append(topicIDs, topicID)
	}

	// 同步每个主题的点踩数
	for _, topicID := range topicIDs {
		likeKey := "topic_likes:" + topicID
		dislikeKey := "topic_dislikes:" + topicID

		// 获取点赞数
		likeCount, err := c.Redis.Get(context.Background(), likeKey).Int64()
		if err != nil && err != redis.Nil {
			continue
		}

		// 获取点踩数
		dislikeCount, err := c.Redis.Get(context.Background(), dislikeKey).Int64()
		if err != nil && err != redis.Nil {
			continue
		}

		// 更新数据库
		c.DB.Model(&models.Topic{}).Where("id = ?", topicID).
			Updates(map[string]interface{}{
				"like_count":    likeCount,
				"dislike_count": dislikeCount,
			})
	}
}

func (c *ForumController) getLikeCount(topicID string) int64 {
	likeKey := "topic_likes:" + topicID
	likeCount, err := c.Redis.Get(context.Background(), likeKey).Int64()
	if err != nil && err != redis.Nil {
		return 0
	}
	return likeCount
}

func (c *ForumController) getDislikeCount(topicID string) int64 {
	dislikeKey := "topic_dislikes:" + topicID
	dislikeCount, err := c.Redis.Get(context.Background(), dislikeKey).Int64()
	if err != nil && err != redis.Nil {
		return 0
	}
	return dislikeCount
}

func (c *ForumController) UnlikeTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// Redis键名
	likeKey := "topic_likes:" + topicID
	userLikeKey := "user_likes:" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// 检查用户是否已点赞
	isLiked, err := c.Redis.SIsMember(ctx, userLikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点赞状态失败"})
		return
	}

	if !isLiked {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未点赞该主题"})
		return
	}

	// 取消点赞
	_, err = c.Redis.SRem(ctx, userLikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	_, err = c.Redis.Decr(ctx, likeKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "取消点赞成功",
		"liked":    false,
		"likes":    c.getLikeCount(topicID),
		"dislikes": c.getDislikeCount(topicID),
	})

	// 立即同步点赞数据
	c.syncAllLikes()
}

func (c *ForumController) UnDislikeTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// Redis键名
	dislikeKey := "topic_dislikes:" + topicID
	userDislikeKey := "user_dislikes:" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// 检查用户是否已点踩
	isDisliked, err := c.Redis.SIsMember(ctx, userDislikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点踩状态失败"})
		return
	}

	if !isDisliked {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未点踩该主题"})
		return
	}

	// 取消点踩
	_, err = c.Redis.SRem(ctx, userDislikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点踩失败"})
		return
	}

	_, err = c.Redis.Decr(ctx, dislikeKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点踩数失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "取消点踩成功",
		"disliked": false,
		"likes":    c.getLikeCount(topicID),
		"dislikes": c.getDislikeCount(topicID),
	})

	// 立即同步点踩数据
	c.syncAllDislikes()
}

func (c *ForumController) LikeTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// Redis键名
	likeKey := "topic_likes:" + topicID
	dislikeKey := "topic_dislikes:" + topicID
	userLikeKey := "user_likes:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	userDislikeKey := "user_dislikes:" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// 检查用户是否已点踩
	isDisliked, err := c.Redis.SIsMember(ctx, userDislikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点踩状态失败"})
		return
	}

	// 如果已点踩，先取消点踩
	if isDisliked {
		// 调试日志：打印当前Redis键值
		fmt.Printf("[DEBUG] Before remove dislike - userDislikeKey: %s, members: %v\n", userDislikeKey, c.Redis.SMembers(ctx, userDislikeKey).Val())
		fmt.Printf("[DEBUG] Before remove dislike - dislikeKey: %s, value: %v\n", dislikeKey, c.Redis.Get(ctx, dislikeKey).Val())

		_, err = c.Redis.SRem(ctx, userDislikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点踩失败"})
			return
		}
		_, err = c.Redis.Decr(ctx, dislikeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点踩数失败"})
			return
		}

		// 调试日志：打印操作后的Redis键值
		fmt.Printf("[DEBUG] After remove dislike - userDislikeKey: %s, members: %v\n", userDislikeKey, c.Redis.SMembers(ctx, userDislikeKey).Val())
		fmt.Printf("[DEBUG] After remove dislike - dislikeKey: %s, value: %v\n", dislikeKey, c.Redis.Get(ctx, dislikeKey).Val())
	}

	// 检查用户是否已点赞
	isLiked, err := c.Redis.SIsMember(ctx, userLikeKey, topicID).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查点赞状态失败"})
		return
	}

	if isLiked {
		// 取消点赞
		_, err = c.Redis.SRem(ctx, userLikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
			return
		}

		_, err = c.Redis.Decr(ctx, likeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "取消点赞成功",
			"liked":    false,
			"disliked": false,
			"likes":    c.getLikeCount(topicID),
			"dislikes": c.getDislikeCount(topicID),
		})
	} else {
		// 点赞
		// 调试日志：打印当前Redis键值
		fmt.Printf("[DEBUG] Before add like - userLikeKey: %s, members: %v\n", userLikeKey, c.Redis.SMembers(ctx, userLikeKey).Val())
		fmt.Printf("[DEBUG] Before add like - likeKey: %s, value: %v\n", likeKey, c.Redis.Get(ctx, likeKey).Val())

		_, err = c.Redis.SAdd(ctx, userLikeKey, topicID).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
			return
		}

		_, err = c.Redis.Incr(ctx, likeKey).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
			return
		}

		// 调试日志：打印操作后的Redis键值
		fmt.Printf("[DEBUG] After add like - userLikeKey: %s, members: %v\n", userLikeKey, c.Redis.SMembers(ctx, userLikeKey).Val())
		fmt.Printf("[DEBUG] After add like - likeKey: %s, value: %v\n", likeKey, c.Redis.Get(ctx, likeKey).Val())

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "点赞成功",
			"liked":    true,
			"disliked": false,
			"likes":    c.getLikeCount(topicID),
			"dislikes": c.getDislikeCount(topicID),
		})
	}

	// 立即同步点赞数据
	c.syncAllLikes()
}

// GetTopicLikes 获取主题点赞数
func (c *ForumController) GetTopicLikes(ctx *gin.Context) {
	topicID := ctx.Param("id")
	likeKey := "topic_likes:" + topicID
	dislikeKey := "topic_dislikes:" + topicID

	// 获取点赞数
	likeCount, err := c.Redis.Get(ctx, likeKey).Int64()
	if err == redis.Nil {
		likeCount = 0
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取点赞数失败"})
		return
	}

	// 获取点踩数
	dislikeCount, err := c.Redis.Get(ctx, dislikeKey).Int64()
	if err == redis.Nil {
		dislikeCount = 0
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取点踩数失败"})
		return
	}

	// 获取用户点赞状态
	userLiked := false
	userDisliked := false
	if userID, exists := ctx.Get("userID"); exists {
		userLikeKey := "user_likes:" + strconv.FormatUint(uint64(userID.(uint)), 10)
		userDislikeKey := "user_dislikes:" + strconv.FormatUint(uint64(userID.(uint)), 10)

		if isLiked, _ := c.Redis.SIsMember(ctx, userLikeKey, topicID).Result(); isLiked {
			userLiked = true
		}
		if isDisliked, _ := c.Redis.SIsMember(ctx, userDislikeKey, topicID).Result(); isDisliked {
			userDisliked = true
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes":        likeCount,
		"dislikes":     dislikeCount,
		"userLiked":    userLiked,
		"userDisliked": userDisliked,
	})
}

// GetTopicById 获取主题详情
func (c *ForumController) GetTopicById(ctx *gin.Context) {
	id := ctx.Param("id")

	var topic models.Topic
	result := c.DB.Preload("User").Preload("Category").First(&topic, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 增加浏览次数
	c.DB.Model(&topic).Update("view_count", gorm.Expr("view_count + 1"))

	// 获取回复
	var replies []models.Reply
	c.DB.Where("topic_id = ?", id).Preload("User").Order("created_at ASC").Find(&replies)

	ctx.JSON(http.StatusOK, gin.H{
		"topic":   topic,
		"replies": replies,
	})
}

// CreateTopic 创建主题
func (c *ForumController) CreateTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if result := c.DB.First(&category, input.CategoryID); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类不存在"})
		return
	}

	// 创建主题
	topic := models.Topic{
		Title:      input.Title,
		Content:    input.Content,
		CategoryID: input.CategoryID,
		UserID:     userID.(uint),
		ViewCount:  0,
		ReplyCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if result := c.DB.Create(&topic); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建主题失败"})
		return
	}

	// 添加积分记录（发帖奖励积分）
	pointRecord := models.PointRecord{
		UserID:      userID.(uint),
		Points:      5, // 发帖奖励5积分
		Type:        "post",
		Description: "发表主题奖励",
		CreatedAt:   time.Now(),
	}

	c.DB.Create(&pointRecord)

	// 更新用户积分
	c.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", 5))

	// 返回创建的主题
	c.DB.Preload("User").Preload("Category").First(&topic, topic.ID)
	ctx.JSON(http.StatusCreated, topic)
}

// UpdateTopic 更新主题
func (c *ForumController) UpdateTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	id := ctx.Param("id")

	// 查询主题
	var topic models.Topic
	result := c.DB.First(&topic, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 检查是否是主题的作者
	if topic.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权修改此主题"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新主题
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Title != "" {
		updates["title"] = input.Title
	}

	if input.Content != "" {
		updates["content"] = input.Content
	}

	if input.CategoryID != 0 {
		// 检查分类是否存在
		var category models.Category
		if result := c.DB.First(&category, input.CategoryID); result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "分类不存在"})
			return
		}
		updates["category_id"] = input.CategoryID
	}

	// 保存更新
	if result := c.DB.Model(&topic).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新主题失败"})
		return
	}

	// 重新查询主题以获取最新信息
	c.DB.Preload("User").Preload("Category").First(&topic, id)

	// 返回更新后的主题
	ctx.JSON(http.StatusOK, topic)
}

// DeleteTopic 删除主题
func (c *ForumController) DeleteTopic(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	id := ctx.Param("id")

	// 查询主题
	var topic models.Topic
	result := c.DB.First(&topic, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 检查是否是主题的作者或管理员
	var user models.User
	c.DB.First(&user, userID)

	if topic.UserID != userID.(uint) && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此主题"})
		return
	}

	// 删除相关回复
	c.DB.Where("topic_id = ?", id).Delete(&models.Reply{})

	// 删除主题
	if result := c.DB.Delete(&topic); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除主题失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "主题已删除"})
}

// CreateReply 创建回复
func (c *ForumController) CreateReply(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取主题ID
	topicID := ctx.Param("id")

	// 检查主题是否存在
	var topic models.Topic
	result := c.DB.First(&topic, topicID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	// 绑定请求数据
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建回复
	reply := models.Reply{
		Content:   input.Content,
		UserID:    userID.(uint),
		TopicID:   topic.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := c.DB.Create(&reply); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建回复失败"})
		return
	}

	// 更新主题回复数
	c.DB.Model(&topic).Update("reply_count", gorm.Expr("reply_count + 1"))

	// 添加积分记录（回复奖励积分）
	pointRecord := models.PointRecord{
		UserID:      userID.(uint),
		Points:      2, // 回复奖励2积分
		Type:        "reply",
		Description: "回复主题奖励",
		CreatedAt:   time.Now(),
	}

	c.DB.Create(&pointRecord)

	// 更新用户积分
	c.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", 2))

	// 返回创建的回复
	c.DB.Preload("User").First(&reply, reply.ID)
	ctx.JSON(http.StatusCreated, reply)
}

// UpdateReply 更新回复
func (c *ForumController) UpdateReply(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取回复ID
	id := ctx.Param("id")

	// 查询回复
	var reply models.Reply
	result := c.DB.First(&reply, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "回复不存在"})
		return
	}

	// 检查是否是回复的作者
	if reply.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权修改此回复"})
		return
	}

	// 绑定请求数据
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新回复
	updates := map[string]interface{}{
		"content":    input.Content,
		"updated_at": time.Now(),
	}

	// 保存更新
	if result := c.DB.Model(&reply).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新回复失败"})
		return
	}

	// 重新查询回复以获取最新信息
	c.DB.Preload("User").First(&reply, id)

	// 返回更新后的回复
	ctx.JSON(http.StatusOK, reply)
}

// DeleteReply 删除回复
func (c *ForumController) DeleteReply(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取回复ID
	id := ctx.Param("id")

	// 查询回复
	var reply models.Reply
	result := c.DB.First(&reply, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "回复不存在"})
		return
	}

	// 检查是否是回复的作者或管理员
	var user models.User
	c.DB.First(&user, userID)

	if reply.UserID != userID.(uint) && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此回复"})
		return
	}

	// 获取主题信息，用于更新回复数
	var topic models.Topic
	c.DB.First(&topic, reply.TopicID)

	// 删除回复
	if result := c.DB.Delete(&reply); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除回复失败"})
		return
	}

	// 更新主题回复数
	c.DB.Model(&topic).Update("reply_count", gorm.Expr("reply_count - 1"))

	ctx.JSON(http.StatusOK, gin.H{"message": "回复已删除"})
}

// GetUserFavorites 获取用户收藏的帖子列表
func (c *ForumController) GetUserFavorites(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取分页和排序参数
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sort := ctx.DefaultQuery("sort", "created_at desc")

	// 转换分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// 计算偏移量
	offset := (page - 1) * limit

	// 查询用户收藏的帖子
	var favorites []models.UserTopicFavorite
	query := c.DB.Preload("Topic").Preload("Topic.User").Preload("Topic.Category")
	query = query.Where("user_id = ?", userID).Order(sort).Offset(offset).Limit(limit)

	if err := query.Find(&favorites).Error; err != nil {
		log.Printf("获取用户收藏失败: %v, 用户ID: %v, 查询参数: page=%d, limit=%d, sort=%s", err, userID, page, limit, sort)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败", "details": "服务器内部错误"})
		return
	}

	// 获取总数
	var total int64
	c.DB.Model(&models.UserTopicFavorite{}).Where("user_id = ?", userID).Count(&total)

	// 格式化返回数据
	var topics []gin.H
	for _, favorite := range favorites {
		// 获取主题信息
		var topic models.Topic
		if err := c.DB.Preload("User").Preload("Category").First(&topic, favorite.TopicID).Error; err != nil {
			continue // 跳过已删除的主题
		}

		topics = append(topics, gin.H{
			"id":          topic.ID,
			"title":       topic.Title,
			"content":     topic.Content,
			"category":    topic.Category,
			"user":        topic.User,
			"view_count":  topic.ViewCount,
			"reply_count": topic.ReplyCount,
			"created_at":  favorite.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": topics,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"totalPage": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}
