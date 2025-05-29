package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"g/front/backend/models"
)

// AdminController 管理员控制器
type AdminController struct {
	DB *gorm.DB
}

// NewAdminController 创建管理员控制器实例
func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

// GetPendingResources 获取待审核资源列表
func (c *AdminController) GetPendingResources(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 查询待审核资源
	var resources []models.Resource
	var total int64

	c.DB.Model(&models.Resource{}).Where("status = ?", "pending").Count(&total)
	c.DB.Where("status = ?", "pending").
		Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resources)

	// 返回资源列表
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// ReviewResource 审核资源
func (c *AdminController) ReviewResource(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取资源ID
	id := ctx.Param("id")

	// 绑定请求数据
	var input struct {
		Status  string `json:"status" binding:"required,oneof=approved rejected"`
		Message string `json:"message"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询资源
	var resource models.Resource
	result := c.DB.First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 更新资源状态
	updates := map[string]interface{}{
		"status":     input.Status,
		"updated_at": time.Now(),
	}

	if result := c.DB.Model(&resource).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新资源状态失败"})
		return
	}

	// 如果审核通过，奖励用户积分
	if input.Status == "approved" {
		// 添加积分记录
		pointRecord := models.PointRecord{
			UserID:      resource.UserID,
			Points:      20, // 资源审核通过奖励20积分
			Type:        "approve",
			ResourceID:  &resource.ID,
			Description: "资源审核通过奖励",
			CreatedAt:   time.Now(),
		}

		c.DB.Create(&pointRecord)

		// 更新用户积分
		c.DB.Model(&models.User{}).Where("id = ?", resource.UserID).Update("points", gorm.Expr("points + ?", 20))
	}

	// 返回更新后的资源
	c.DB.Preload("User").Preload("Category").First(&resource, id)
	ctx.JSON(http.StatusOK, resource)
}

// GetUserStats 获取用户统计信息
func (c *AdminController) GetUserStats(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 统计用户数量
	var userCount int64
	c.DB.Model(&models.User{}).Count(&userCount)

	// 统计资源数量
	var resourceCount int64
	c.DB.Model(&models.Resource{}).Count(&resourceCount)

	// 统计待审核资源数量
	var pendingCount int64
	c.DB.Model(&models.Resource{}).Where("status = ?", "pending").Count(&pendingCount)

	// 统计今日新增用户
	var todayUserCount int64
	c.DB.Model(&models.User{}).Where("created_at >= ?", time.Now().Format("2006-01-02")).Count(&todayUserCount)

	// 返回统计信息
	ctx.JSON(http.StatusOK, gin.H{
		"user_count":       userCount,
		"resource_count":   resourceCount,
		"pending_count":    pendingCount,
		"today_user_count": todayUserCount,
	})
}

// GetUsers 获取用户列表
func (c *AdminController) GetUsers(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	sort := ctx.DefaultQuery("sort", "created_at:desc")
	search := ctx.Query("search")
	role := ctx.Query("role")

	// 构建查询
	query := c.DB.Model(&models.User{})

	// 搜索条件
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 排序
	if sort != "" {
		parts := strings.Split(sort, ":")
		if len(parts) == 2 {
			field := parts[0]
			order := parts[1]
			if order == "desc" {
				query = query.Order(field + " DESC")
			} else {
				query = query.Order(field)
			}
		}
	}

	// 查询用户总数
	var total int64
	query.Count(&total)

	// 查询用户列表
	var users []models.User
	query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)

	// 返回用户列表
	ctx.JSON(http.StatusOK, gin.H{
		"users":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteUser 删除用户
func (c *AdminController) DeleteUser(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取用户ID
	id := ctx.Param("id")

	// 不能删除自己
	if id == strconv.FormatUint(uint64(admin.ID), 10) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账号"})
		return
	}

	// 删除用户
	result := c.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// UpdateUserRole 更新用户角色
func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取用户ID
	id := ctx.Param("id")

	// 不能修改自己的角色
	if id == strconv.FormatUint(uint64(admin.ID), 10) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能修改自己的角色"})
		return
	}

	// 绑定请求数据
	var input struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户角色
	result := c.DB.Model(&models.User{}).Where("id = ?", id).Update("role", input.Role)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户角色失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户角色已更新"})
}

// GetResources 获取资源列表
func (c *AdminController) GetResources(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	sort := ctx.DefaultQuery("sort", "created_at:desc")
	search := ctx.Query("search")
	status := ctx.Query("status")
	categoryID := ctx.Query("category_id")

	// 构建查询
	query := c.DB.Model(&models.Resource{})

	// 搜索条件
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 分类筛选
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 排序
	if sort != "" {
		parts := strings.Split(sort, ":")
		if len(parts) == 2 {
			field := parts[0]
			order := parts[1]
			if order == "desc" {
				query = query.Order(field + " DESC")
			} else {
				query = query.Order(field)
			}
		}
	}

	// 查询资源总数
	var total int64
	query.Count(&total)

	// 查询资源列表
	var resources []models.Resource
	query.Preload("User").Preload("Category").Limit(pageSize).Offset((page - 1) * pageSize).Find(&resources)

	// 返回资源列表
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// DeleteResource 删除资源
func (c *AdminController) DeleteResource(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取资源ID
	id := ctx.Param("id")

	// 删除资源
	result := c.DB.Delete(&models.Resource{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除资源失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "资源已删除"})
}

// GetTopics 获取论坛话题列表
func (c *AdminController) GetTopics(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	sort := ctx.DefaultQuery("sort", "created_at:desc")
	search := ctx.Query("search")
	categoryID := ctx.Query("category_id")

	// 构建查询
	query := c.DB.Model(&models.Topic{})

	// 搜索条件
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 分类筛选
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 排序
	if sort != "" {
		parts := strings.Split(sort, ":")
		if len(parts) == 2 {
			field := parts[0]
			order := parts[1]
			if order == "desc" {
				query = query.Order(field + " DESC")
			} else {
				query = query.Order(field)
			}
		}
	}

	// 查询话题总数
	var total int64
	query.Count(&total)

	// 查询话题列表
	var topics []models.Topic
	query.Preload("User").Preload("Category").Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(3)
	}).Preload("Replies.User").Limit(pageSize).Offset((page - 1) * pageSize).Find(&topics)

	// 返回话题列表
	ctx.JSON(http.StatusOK, gin.H{
		"topics":   topics,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteTopic 删除论坛话题
func (c *AdminController) DeleteTopic(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取话题ID
	id := ctx.Param("id")

	// 删除话题
	result := c.DB.Delete(&models.Topic{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除话题失败"})
		return
	}

	// 删除相关回复
	c.DB.Where("topic_id = ?", id).Delete(&models.Reply{})

	ctx.JSON(http.StatusOK, gin.H{"message": "话题已删除"})
}

// GetPointsRecords 获取积分记录列表
func (c *AdminController) GetPointsRecords(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	sort := ctx.DefaultQuery("sort", "created_at:desc")
	search := ctx.Query("search")
	pointType := ctx.Query("type")

	// 构建查询
	query := c.DB.Model(&models.PointRecord{})

	// 搜索条件（按用户名搜索）
	if search != "" {
		var userIds []uint
		c.DB.Model(&models.User{}).Where("username LIKE ?", "%"+search+"%").Pluck("id", &userIds)
		if len(userIds) > 0 {
			query = query.Where("user_id IN ?", userIds)
		} else {
			// 如果没有找到匹配的用户，返回空结果
			ctx.JSON(http.StatusOK, gin.H{
				"records":  []models.PointRecord{},
				"total":    0,
				"page":     page,
				"pageSize": pageSize,
			})
			return
		}
	}

	// 类型筛选
	if pointType != "" {
		query = query.Where("type = ?", pointType)
	}

	// 排序
	if sort != "" {
		parts := strings.Split(sort, ":")
		if len(parts) == 2 {
			field := parts[0]
			order := parts[1]
			if order == "desc" {
				query = query.Order(field + " DESC")
			} else {
				query = query.Order(field)
			}
		}
	}

	// 查询积分记录总数
	var total int64
	query.Count(&total)

	// 查询积分记录列表
	var records []models.PointRecord
	query.Preload("User").Limit(pageSize).Offset((page - 1) * pageSize).Find(&records)

	// 返回积分记录列表
	ctx.JSON(http.StatusOK, gin.H{
		"records":  records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// AdjustPoints 调整用户积分
func (c *AdminController) AdjustPoints(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 绑定请求数据
	var input struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Points      int    `json:"points" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户是否存在
	var user models.User
	result := c.DB.First(&user, input.UserID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 添加积分记录
	pointRecord := models.PointRecord{
		UserID:      input.UserID,
		Points:      input.Points,
		Type:        "admin",
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	c.DB.Create(&pointRecord)

	// 更新用户积分
	c.DB.Model(&models.User{}).Where("id = ?", input.UserID).Update("points", gorm.Expr("points + ?", input.Points))

	ctx.JSON(http.StatusOK, gin.H{"message": "积分已调整"})
}

// GetResourceStats 获取资源统计信息
func (c *AdminController) GetResourceStats(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取资源总数
	var totalResources int64
	c.DB.Model(&models.Resource{}).Count(&totalResources)

	// 获取各类型资源数量
	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	c.DB.Model(&models.Resource{}).Select("type, count(*) as count").Group("type").Scan(&typeStats)

	// 获取各状态资源数量
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	c.DB.Model(&models.Resource{}).Select("status, count(*) as count").Group("status").Scan(&statusStats)

	// 获取最近7天的资源上传数量
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 计算7天前的日期
	sevenDaysAgo := time.Now().AddDate(0, 0, -6)

	// 查询每天的资源上传数量
	for i := 0; i < 7; i++ {
		date := sevenDaysAgo.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		nextDate := date.AddDate(0, 0, 1)

		var count int64
		c.DB.Model(&models.Resource{}).Where("created_at >= ? AND created_at < ?", date, nextDate).Count(&count)

		dailyStats = append(dailyStats, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{
			Date:  dateStr,
			Count: count,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":        totalResources,
		"type_stats":   typeStats,
		"status_stats": statusStats,
		"daily_stats":  dailyStats,
	})
}

// GetForumStats 获取论坛统计信息
func (c *AdminController) GetForumStats(ctx *gin.Context) {
	// 验证是否为管理员
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var admin models.User
	c.DB.First(&admin, userID)
	if admin.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取话题总数
	var totalTopics int64
	c.DB.Model(&models.Topic{}).Count(&totalTopics)

	// 获取回复总数
	var totalReplies int64
	c.DB.Model(&models.Reply{}).Count(&totalReplies)

	// 获取各分类话题数量
	var categoryStats []struct {
		CategoryID   uint   `json:"category_id"`
		CategoryName string `json:"category_name"`
		Count        int64  `json:"count"`
	}

	c.DB.Model(&models.Topic{}).
		Select("topics.category_id, categories.name as category_name, count(*) as count").
		Joins("left join categories on topics.category_id = categories.id").
		Group("topics.category_id").
		Scan(&categoryStats)

	// 获取最近7天的话题发布数量
	var dailyTopicStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 获取最近7天的回复发布数量
	var dailyReplyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 计算7天前的日期
	sevenDaysAgo := time.Now().AddDate(0, 0, -6)

	// 查询每天的话题和回复数量
	for i := 0; i < 7; i++ {
		date := sevenDaysAgo.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		nextDate := date.AddDate(0, 0, 1)

		// 话题数量
		var topicCount int64
		c.DB.Model(&models.Topic{}).Where("created_at >= ? AND created_at < ?", date, nextDate).Count(&topicCount)

		dailyTopicStats = append(dailyTopicStats, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{
			Date:  dateStr,
			Count: topicCount,
		})

		// 回复数量
		var replyCount int64
		c.DB.Model(&models.Reply{}).Where("created_at >= ? AND created_at < ?", date, nextDate).Count(&replyCount)

		dailyReplyStats = append(dailyReplyStats, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{
			Date:  dateStr,
			Count: replyCount,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_topics":      totalTopics,
		"total_replies":     totalReplies,
		"category_stats":    categoryStats,
		"daily_topic_stats": dailyTopicStats,
		"daily_reply_stats": dailyReplyStats,
	})
}
