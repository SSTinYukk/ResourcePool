package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"gorm.io/gorm"

	"g/front/backend/config"
	"g/front/backend/models"
)

// ChatController AI聊天控制器
type ChatController struct {
	DB       *gorm.DB
	AIConfig config.AIConfig
}

// NewChatController 创建聊天控制器实例
func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{
		DB:       db,
		AIConfig: config.GetAIConfig(),
	}
}

// CreateSession 创建新的聊天会话
func (c *ChatController) CreateSession(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title string `json:"title" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建新会话
	session := models.ChatSession{
		UserID:    userID.(uint),
		Title:     input.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := c.DB.Create(&session)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 返回创建的会话
	ctx.JSON(http.StatusCreated, session)
}

// GetSessions 获取用户的所有聊天会话
func (c *ChatController) GetSessions(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 查询用户的会话
	var sessions []models.ChatSession
	var total int64

	c.DB.Model(&models.ChatSession{}).Where("user_id = ?", userID).Count(&total)
	c.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&sessions)

	// 返回会话列表
	ctx.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetSessionMessages 获取会话的消息历史
func (c *ChatController) GetSessionMessages(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取会话ID
	sessionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 验证会话所有权
	var session models.ChatSession
	result := c.DB.First(&session, sessionID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	if session.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "50"))

	// 查询会话消息
	var messages []models.ChatMessage
	var total int64

	c.DB.Model(&models.ChatMessage{}).Where("session_id = ?", sessionID).Count(&total)
	c.DB.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&messages)

	// 返回消息列表
	ctx.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteSession 删除聊天会话
func (c *ChatController) DeleteSession(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取会话ID
	sessionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 验证会话所有权
	var session models.ChatSession
	result := c.DB.First(&session, sessionID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	if session.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此会话"})
		return
	}

	// 删除会话及其消息（依赖于外键级联删除）
	c.DB.Delete(&session)

	// 返回成功消息
	ctx.JSON(http.StatusOK, gin.H{"message": "会话已删除"})
}

// SendMessage 发送消息并获取AI回复
func (c *ChatController) SendMessage(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 绑定请求数据
	var input struct {
		SessionID uint   `json:"session_id" binding:"required,gt=0"`
		Content   string `json:"content" binding:"required,min=1,max=1000"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 验证会话所有权
	var session models.ChatSession
	result := c.DB.First(&session, input.SessionID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	if session.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 创建用户消息
	userMessage := models.ChatMessage{
		SessionID: input.SessionID,
		UserID:    userID.(uint),
		Role:      "user",
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	c.DB.Create(&userMessage)

	// 获取会话历史消息（最近10条）
	var history []models.ChatMessage
	c.DB.Where("session_id = ?", input.SessionID).
		Order("created_at DESC").
		Limit(10).
		Find(&history)

	// 构建请求体
	messages := []map[string]string{}
	// 按时间正序添加历史消息
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, map[string]string{
			"role":    history[i].Role,
			"content": history[i].Content,
		})
	}

	// 调用AI API获取回复
	response, err := c.callAIAPI(messages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务调用失败: " + err.Error()})
		return
	}

	// 保存AI回复
	aiMessage := models.ChatMessage{
		SessionID: input.SessionID,
		UserID:    userID.(uint),
		Role:      "assistant",
		Content:   response,
		CreatedAt: time.Now(),
	}

	c.DB.Create(&aiMessage)

	// 更新会话最后更新时间
	c.DB.Model(&session).Update("updated_at", time.Now())

	// 返回AI回复
	ctx.JSON(http.StatusOK, aiMessage)
}

// callAIAPI 调用火山引擎AI API
func (c *ChatController) callAIAPI(messages []map[string]string) (string, error) {
	// 创建火山引擎客户端
	client := arkruntime.NewClientWithApiKey(c.AIConfig.APIKey)
	ctx := context.Background()

	// 转换消息格式
	var sdkMessages []*model.ChatCompletionMessage
	for _, msg := range messages {
		sdkMessages = append(sdkMessages, &model.ChatCompletionMessage{
			Role: *volcengine.String(msg["role"]),
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(msg["content"]),
			},
		})
	}

	// 构建请求
	req := model.ChatCompletionRequest{
		Model:    c.AIConfig.ModelID,
		Messages: sdkMessages,
	}

	// 调用API
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("API调用失败: %v", err)
	}

	// 获取响应内容
	if resp.Choices == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("无效的API响应格式: 缺少choices字段")
	}

	// if resp.Choices[0].Message == nil || resp.Choices[0].Message.Content == nil || resp.Choices[0].Message.Content.StringValue == nil {
	// 	return "", fmt.Errorf("无效的API响应格式: 缺少message或content字段")
	// }

	return *resp.Choices[0].Message.Content.StringValue, nil
}
