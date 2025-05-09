package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"g/front/backend/middleware"
	"g/front/backend/models"
)

// UserController 用户控制器
type UserController struct {
	DB *gorm.DB
}

// NewUserController 创建用户控制器实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if result := c.DB.Where("username = ?", input.Username).First(&existingUser); result.Error == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名已被使用"})
		return
	}

	// 检查邮箱是否已存在
	if result := c.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Points:    0,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := c.DB.Create(&user); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"points":   user.Points,
			"role":     user.Role,
		},
		"token": token,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"` // 用户名或邮箱
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	result := c.DB.Where("username = ? OR email = ?", input.Username, input.Username).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"points":   user.Points,
			"role":     user.Role,
		},
		"token": token,
	})
}

// GetProfile 获取用户资料
func (c *UserController) GetProfile(ctx *gin.Context) {
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

	// 返回用户资料
	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"points":   user.Points,
		"role":     user.Role,
	})
}

// UpdateProfile 更新用户资料
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var input struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.User
	result := c.DB.First(&user, userID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Email != "" && input.Email != user.Email {
		// 检查邮箱是否已被使用
		var existingUser models.User
		if result := c.DB.Where("email = ? AND id != ?", input.Email, userID).First(&existingUser); result.Error == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
			return
		}
		updates["email"] = input.Email
	}

	if input.Avatar != "" {
		updates["avatar"] = input.Avatar
	}

	// 保存更新
	if result := c.DB.Model(&user).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户资料失败"})
		return
	}

	// 重新查询用户以获取最新信息
	c.DB.First(&user, userID)

	// 返回更新后的用户资料
	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"points":   user.Points,
		"role":     user.Role,
	})
}

// GetPoints 获取用户积分
func (c *UserController) GetPoints(ctx *gin.Context) {
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

	// 返回用户积分
	ctx.JSON(http.StatusOK, gin.H{
		"points": user.Points,
	})
}

// GetPointHistory 获取积分历史
func (c *UserController) GetPointHistory(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	// 将字符串类型的 page 和 pageSize 转换为 int 类型
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分页参数 page 格式错误"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "分页参数 pageSize 格式错误"})
		return
	}

	// 查询积分历史记录
	var records []models.PointRecord
	var total int64

	c.DB.Model(&models.PointRecord{}).Where("user_id = ?", userID).Count(&total)
	c.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&records)

	// 返回积分历史
	ctx.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   total,
	})
}
