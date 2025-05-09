package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"g/front/backend/models"
)

// ForumController 论坛控制器
type ForumController struct {
	DB *gorm.DB
}

// NewForumController 创建论坛控制器实例
func NewForumController(db *gorm.DB) *ForumController {
	return &ForumController{DB: db}
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
