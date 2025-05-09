package controllers

import (
	"net/http"
	"strconv"
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
