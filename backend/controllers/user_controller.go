package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"g/front/backend/middleware"
	"g/front/backend/models"
)

// UserController 用户控制器
type UserController struct {
	DB          *gorm.DB
	MinioClient *minio.Client
}

// NewUserController 创建用户控制器实例
func NewUserController(db *gorm.DB, minioClient *minio.Client) *UserController {
	return &UserController{DB: db, MinioClient: minioClient}
}

// GetUserProfile 获取用户资料
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := c.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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

// RefreshToken 刷新JWT令牌
func (c *UserController) RefreshToken(ctx *gin.Context) {
	// 从上下文中获取用户ID
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

	// 生成新的JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回新的令牌
	ctx.JSON(http.StatusOK, gin.H{
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

// UploadAvatar 上传用户头像
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户当前头像
	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 如果用户已有头像，从MinIO中删除旧头像
	if user.Avatar != "" {
		// 从URL中提取文件名
		oldFileName := strings.TrimPrefix(user.Avatar, "http://47.121.210.209:9000/avatars/")
		err := c.MinioClient.RemoveObject(ctx, "avatars", oldFileName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("删除旧头像失败: %v", err)
		}
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请选择头像文件"})
		return
	}
	defer file.Close()

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[header.Header.Get("Content-Type")] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "仅支持JPEG、PNG或GIF格式的图片"})
		return
	}

	// 验证文件大小 (限制5MB)
	if header.Size > 5<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "头像文件大小不能超过5MB"})
		return
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("avatars/%d-%s", time.Now().UnixNano(), header.Filename)

	// 检查并创建avatars存储桶
	bucketExists, err := c.MinioClient.BucketExists(ctx, "avatars")
	if err != nil {
		log.Printf("检查存储桶失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误", "details": err.Error()})
		return
	}
	if !bucketExists {
		err = c.MinioClient.MakeBucket(ctx, "avatars", minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("创建存储桶失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误", "details": err.Error()})
			return
		}
	}

	// 上传到MinIO
	_, err = c.MinioClient.PutObject(ctx, "avatars", fileName, file, header.Size, minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")})
	if err != nil {
		log.Printf("MinIO上传失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "上传头像失败", "details": err.Error()})
		return
	}

	// 更新用户头像URL
	avatarURL := fmt.Sprintf("http://47.121.210.209:9000/avatars/%s", fileName)
	if result := c.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL); result.Error != nil {
		log.Printf("数据库更新失败: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败", "details": result.Error.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"avatar":  avatarURL,
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
