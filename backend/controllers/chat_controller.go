package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"g/front/backend/models"
)

// ChatController AI聊天控制器
type ChatController struct {
	DB *gorm.DB
}

// NewChatController 创建聊天控制器实例
func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{DB: db}
}

// GetChatHistories 获取聊天历史列表
func (c *ChatController) GetChatHistories(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询聊天历史
	var histories []models.ChatHistory
	c.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&histories)

	// 返回结果
	ctx.JSON(http.StatusOK, histories)
}

// GetChatHistoryById 获取聊天历史详情
func (c *ChatController) GetChatHistoryById(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取聊天历史ID
	id := ctx.Param("id")

	// 查询聊天历史
	var history models.ChatHistory
	result := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&history)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "聊天历史不存在"})
		return
	}

	// 查询聊天消息
	var messages []models.ChatMessage
	c.DB.Where("chat_history_id = ?", id).
		Order("created_at ASC").
		Find(&messages)

	// 返回结果
	history.Messages = messages
	ctx.JSON(http.StatusOK, history)
}

// CreateChatHistory 创建新的聊天历史
func (c *ChatController) CreateChatHistory(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title string `json:"title"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建聊天历史
	title := input.Title
	if title == "" {
		title = "新对话"
	}

	chatHistory := models.ChatHistory{
		UserID:    userID.(uint),
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := c.DB.Create(&chatHistory); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建聊天历史失败"})
		return
	}

	// 返回创建的聊天历史
	ctx.JSON(http.StatusCreated, chatHistory)
}

// UpdateChatHistory 更新聊天历史
func (c *ChatController) UpdateChatHistory(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取聊天历史ID
	id := ctx.Param("id")

	// 查询聊天历史
	var history models.ChatHistory
	result := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&history)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "聊天历史不存在"})
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

	// 更新聊天历史
	updates := map[string]interface{}{
		"title":      input.Title,
		"updated_at": time.Now(),
	}

	if result := c.DB.Model(&history).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新聊天历史失败"})
		return
	}

	// 返回更新后的聊天历史
	c.DB.First(&history, id)
	ctx.JSON(http.StatusOK, history)
}

// DeleteChatHistory 删除聊天历史
func (c *ChatController) DeleteChatHistory(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取聊天历史ID
	id := ctx.Param("id")

	// 查询聊天历史
	var history models.ChatHistory
	result := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&history)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "聊天历史不存在"})
		return
	}

	// 删除相关聊天消息
	c.DB.Where("chat_history_id = ?", id).Delete(&models.ChatMessage{})

	// 删除聊天历史
	if result := c.DB.Delete(&history); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除聊天历史失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "聊天历史已删除"})
}

// SendMessage 发送消息
func (c *ChatController) SendMessage(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取聊天历史ID
	id := ctx.Param("id")

	// 查询聊天历史
	var history models.ChatHistory
	result := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&history)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "聊天历史不存在"})
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

	// 创建用户消息
	userMessage := models.ChatMessage{
		ChatHistoryID: history.ID,
		Sender:        "user",
		Content:       input.Content,
		CreatedAt:     time.Now(),
	}

	if result := c.DB.Create(&userMessage); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "发送消息失败"})
		return
	}

	// 更新聊天历史标题（如果是第一条消息）
	var messageCount int64
	c.DB.Model(&models.ChatMessage{}).Where("chat_history_id = ?", history.ID).Count(&messageCount)

	if messageCount == 1 {
		title := input.Content
		if len(title) > 20 {
			title = title[:20] + "..."
		}

		c.DB.Model(&history).Updates(map[string]interface{}{
			"title":      title,
			"updated_at": time.Now(),
		})
	} else {
		c.DB.Model(&history).Update("updated_at", time.Now())
	}

	// 生成AI回复
	aiResponse := c.generateAIResponse(input.Content)

	// 创建AI消息
	aiMessage := models.ChatMessage{
		ChatHistoryID: history.ID,
		Sender:        "ai",
		Content:       aiResponse,
		CreatedAt:     time.Now(),
	}

	if result := c.DB.Create(&aiMessage); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "AI回复失败"})
		return
	}

	// 返回消息
	ctx.JSON(http.StatusOK, gin.H{
		"user_message": userMessage,
		"ai_message":   aiMessage,
	})
}

// generateAIResponse 生成AI回复（模拟）
func (c *ChatController) generateAIResponse(message string) string {
	// 这里是模拟回复，实际项目中应该调用AI服务
	responses := map[string]string{
		"什么是微处理器": "微处理器是一个集成电路芯片，它包含了计算机的中央处理单元(CPU)的功能。它是计算机系统的核心，负责执行指令、处理数据和控制系统的其他部分。微处理器通常包含算术逻辑单元(ALU)、控制单元、寄存器组和内部总线等组件。",

		"8086CPU的寻址方式": "8086CPU的主要寻址方式包括：<br>1. <b>立即寻址</b>：操作数直接包含在指令中<br>2. <b>寄存器寻址</b>：操作数在CPU寄存器中<br>3. <b>直接寻址</b>：指令中给出操作数的有效地址<br>4. <b>寄存器间接寻址</b>：寄存器中存放操作数的有效地址<br>5. <b>基址寻址</b>：有效地址=基址寄存器内容+位移量<br>6. <b>变址寻址</b>：有效地址=变址寄存器内容+位移量<br>7. <b>基址变址寻址</b>：有效地址=基址寄存器内容+变址寄存器内容+位移量",

		"中断向量表": "中断向量表是存储在内存中的一个表，包含了各种中断服务程序的入口地址。当中断发生时，CPU会根据中断类型查找中断向量表中对应的入口地址，然后跳转到该地址执行中断服务程序。<br><br>在8086/8088中，中断向量表位于内存的最低1KB(0000H-03FFH)，每个中断向量占4个字节，可以存储256个中断向量。每个向量的低2字节是偏移地址，高2字节是段地址。",

		"微机系统的基本组成": "微机系统的基本组成包括：<br>1. <b>中央处理器(CPU)</b>：执行指令、处理数据<br>2. <b>存储器</b>：包括内存(RAM)和只读存储器(ROM)<br>3. <b>输入/输出设备</b>：键盘、显示器、打印机等<br>4. <b>总线系统</b>：连接各个部件，包括数据总线、地址总线和控制总线<br>5. <b>接口电路</b>：连接CPU与外部设备<br><br>这些组件协同工作，实现数据的输入、处理、存储和输出功能。",

		"汇编语言的基本指令": "汇编语言的基本指令类型包括：<br>1. <b>数据传送指令</b>：MOV, PUSH, POP, XCHG等<br>2. <b>算术运算指令</b>：ADD, SUB, MUL, DIV, INC, DEC等<br>3. <b>逻辑运算指令</b>：AND, OR, XOR, NOT等<br>4. <b>移位指令</b>：SHL, SHR, ROL, ROR等<br>5. <b>转移指令</b>：JMP, JZ, JNZ, CALL, RET等<br>6. <b>串操作指令</b>：MOVS, CMPS, SCAS等<br>7. <b>处理器控制指令</b>：CLC, STC, CLI, STI等<br><br>这些指令是汇编语言编程的基础，通过组合这些指令可以实现各种复杂的功能。",
	}

	// 检查是否有预设回复
	for key, response := range responses {
		if containsIgnoreCase(message, key) {
			return response
		}
	}

	// 默认回复
	return "关于\"\" + message + \"\"的问题，我的回答是：<br><br>这是一个关于微机原理的重要概念。在微机系统中，这涉及到处理器架构、指令集和系统总线等核心知识。<br><br>建议您查阅教材的相关章节，或者在论坛中与其他同学讨论这个问题，以获取更详细的解答。<br><br>您还有其他问题吗？"
}

// containsIgnoreCase 忽略大小写检查字符串是否包含子串
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
