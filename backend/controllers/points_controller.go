package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"g/front/backend/models"
)

// PointsController 积分控制器
type PointsController struct {
	DB *gorm.DB
}

// NewPointsController 创建积分控制器实例
func NewPointsController(db *gorm.DB) *PointsController {
	return &PointsController{DB: db}
}

// GetUserPoints 获取用户积分信息
func (c *PointsController) GetUserPoints(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	result := c.DB.First(&user, userID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 返回用户积分信息
	ctx.JSON(http.StatusOK, gin.H{
		"points": user.Points,
	})
}

// GetPointsHistory 获取用户积分历史记录
func (c *PointsController) GetPointsHistory(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 查询积分记录
	var records []models.PointRecord
	var total int64

	c.DB.Model(&models.PointRecord{}).Where("user_id = ?", userID).Count(&total)
	c.DB.Where("user_id = ?", userID).
		Preload("Resource").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&records)

	// 返回积分记录
	ctx.JSON(http.StatusOK, gin.H{
		"records":  records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// AddPoints 添加积分（管理员功能）
func (c *PointsController) AddPoints(ctx *gin.Context) {
	// 从上下文获取用户ID（管理员）
	adminID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 验证是否为管理员
	var admin models.User
	c.DB.First(&admin, adminID)
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

	// 检查用户是否存在
	var user models.User
	if result := c.DB.First(&user, input.UserID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 创建积分记录
	pointRecord := models.PointRecord{
		UserID:      input.UserID,
		Points:      input.Points,
		Type:        "admin",
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	// 开始事务
	tx := c.DB.Begin()

	// 创建积分记录
	if result := tx.Create(&pointRecord); result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建积分记录失败"})
		return
	}

	// 更新用户积分
	if result := tx.Model(&user).Update("points", gorm.Expr("points + ?", input.Points)); result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户积分失败"})
		return
	}

	// 提交事务
	tx.Commit()

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "积分添加成功",
		"pointRecord": pointRecord,
		"newPoints":   user.Points + input.Points,
	})
}

// DeductPoints 扣除积分（内部方法）
func (c *PointsController) DeductPoints(userID uint, points int, recordType string, resourceID *uint, description string) error {
	// 开始事务
	tx := c.DB.Begin()

	// 检查用户积分是否足够
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.Points < points {
		tx.Rollback()
		return fmt.Errorf("积分不足")
	}

	// 创建积分记录
	pointRecord := models.PointRecord{
		UserID:      userID,
		Points:      -points, // 负值表示扣除
		Type:        recordType,
		ResourceID:  resourceID,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := tx.Create(&pointRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新用户积分
	if err := tx.Model(&user).Update("points", gorm.Expr("points - ?", points)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}
