package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"g/front/backend/config"
	"g/front/backend/models"
)

// ResourceController 资源控制器
type ResourceController struct {
	DB          *gorm.DB
	MinioClient *minio.Client
}

// NewResourceController 创建资源控制器实例
func NewResourceController(db *gorm.DB, minioClient *minio.Client) *ResourceController {
	return &ResourceController{DB: db, MinioClient: minioClient}
}

// GetResources 获取资源列表
func (c *ResourceController) GetResources(ctx *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	categoryID := ctx.Query("category")
	sort := ctx.DefaultQuery("sort", "newest")

	// 构建查询
	query := c.DB.Model(&models.Resource{}).Where("status = ?", "approved")

	// 分类过滤
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 排序
	switch sort {
	case "newest":
		query = query.Order("created_at DESC")
	case "popular":
		query = query.Order("download_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 执行查询
	var resources []models.Resource
	var total int64

	query.Count(&total)
	query.Preload("User").Preload("Category").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resources)

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// GetResourceById 获取资源详情
func (c *ResourceController) GetResourceById(ctx *gin.Context) {
	id := ctx.Param("id")

	var resource models.Resource
	result := c.DB.Preload("User").Preload("Category").First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 增加浏览次数逻辑可以在这里添加

	ctx.JSON(http.StatusOK, resource)
}

// GetCategories 获取资源分类
func (c *ResourceController) GetCategories(ctx *gin.Context) {
	var categories []models.Category
	c.DB.Find(&categories)

	ctx.JSON(http.StatusOK, categories)
}

// GetComments 获取资源评论
func (c *ResourceController) GetComments(ctx *gin.Context) {
	// 获取资源ID
	resourceID := ctx.Param("id")

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 查询评论
	var comments []models.Comment
	var total int64

	c.DB.Model(&models.Comment{}).Where("resource_id = ?", resourceID).Count(&total)
	c.DB.Where("resource_id = ?", resourceID).
		Preload("User").
		Order("time DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&comments)

	// 返回评论列表
	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// PostComment 提交资源评论
func (c *ResourceController) PostComment(ctx *gin.Context) {
	// 获取资源ID
	resourceID := ctx.Param("id")

	// 验证用户身份
	log.Printf("请求头信息: %+v", ctx.Request.Header)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}

	// 解析请求数据
	var request struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Content string `json:"content" binding:"required,min=1,max=500"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建评论
	comment := models.Comment{
		ResourceID: resourceID,
		UserID:     userID,
		Rating:     request.Rating,
		Content:    request.Content,
		Time:       time.Now(),
	}

	// 保存到数据库
	if err := c.DB.Create(&comment).Error; err != nil {
		log.Printf("Failed to create comment: %v, request data: %+v", err, request)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "评论提交失败", "details": err.Error()})
		return
	}

	// 返回新创建的评论
	ctx.JSON(http.StatusCreated, gin.H{
		"comment": comment,
	})
}

// UploadResource 上传资源文件
func (c *ResourceController) UploadResource(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传文件",
		})
		return
	}
	defer file.Close()

	// 生成唯一文件名
	fileName := uuid.New().String() + filepath.Ext(header.Filename)
	bucketName := "resources"

	// 确保存储桶存在
	err = c.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// 检查存储桶是否已存在
		exists, errBucketExists := c.MinioClient.BucketExists(context.Background(), bucketName)
		if !exists || errBucketExists != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "无法创建或访问存储桶",
				"error":   errBucketExists.Error(),
			})
			return
		}
	}

	// 开始事务
	tx := c.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "数据库事务启动失败",
			"error":   tx.Error.Error(),
		})
		return
	}

	// 上传到MinIO
	_, err = c.MinioClient.PutObject(
		context.Background(),
		bucketName,
		fileName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")},
	)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "文件上传失败",
			"error":   err.Error(),
		})
		return
	}

	// 从请求中获取资源信息
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	categoryIDStr := ctx.PostForm("category_id")

	// 验证必填字段
	if title == "" || description == "" || categoryIDStr == "" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "标题、描述和分类ID都是必填项",
		})
		return
	}

	// 验证分类ID格式
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "分类ID格式错误",
		})
		return
	}
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权，请先登录",
		})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "用户ID格式错误",
		})
		return
	}

	// 创建资源记录
	resource := models.Resource{
		Title:       title,
		Description: description,
		CategoryID:  uint(categoryID),
		FilePath:    fileName,
		FileSize:    header.Size,
		FileType:    header.Header.Get("Content-Type"),
		Status:      "pending",
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	log.Printf("准备插入资源记录: %+v\n", resource)
	log.Printf("数据库连接状态: %v\n", c.DB.Exec("SELECT 1").Error)
	// 使用Ping方法检查数据库连接状态
	sqlDB, err := c.DB.DB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
	} else {
		if err := sqlDB.Ping(); err != nil {
			log.Printf("数据库连接检查失败: %v", err)
		}
	}
	log.Printf("资源表结构: %+v\n", models.Resource{})
	log.Printf("事务状态: %v\n", tx.Error)
	log.Printf("分类ID验证结果: %d\n", categoryID)
	log.Printf("用户ID验证结果: %d\n", userID)

	// 检查分类是否存在
	var category models.Category
	if result := tx.First(&category, categoryID); result.Error != nil {
		tx.Rollback()
		log.Printf("分类不存在: %v", result.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "分类不存在",
		})
		return
	}

	if result := tx.Create(&resource); result.Error != nil {
		tx.Rollback()
		log.Printf("资源插入失败: %v", result.Error)
		log.Printf("完整错误信息: %+v", result)
		log.Printf("当前资源对象: %+v", resource)
		log.Printf("数据库最后错误: %v", tx.Error)

		// 检查数据库连接状态
		sqlDB, err := tx.DB()
		if err != nil {
			log.Printf("获取数据库连接失败: %v", err)
		} else {
			if err := sqlDB.Ping(); err != nil {
				log.Printf("数据库连接检查失败: %v", err)
			}
		}

		// 检查表结构是否匹配
		var tableSchema string
		tx.Raw("SHOW CREATE TABLE resources").Scan(&tableSchema)
		log.Printf("表结构: %s", tableSchema)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "资源创建失败",
			"error":   result.Error.Error(),
			"details": gin.H{
				"database_status": sqlDB.Stats(),
				"table_schema":    tableSchema,
			},
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("事务提交失败: %v", err)
		log.Printf("事务隔离级别: %v", tx.Statement.DB.Exec("SELECT @@tx_isolation").Error)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "资源创建失败",
			"error":   err.Error(),
			"details": gin.H{
				"transaction_status": "failed",
			},
		})
		return
	}

	// 返回文件URL和资源ID
	endpoint := config.GetEnv("MINIO_ENDPOINT", "47.121.210.209:9000")
	fileURL := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, fileName)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传成功",
		"data": gin.H{
			"url":  fileURL,
			"id":   resource.ID,
			"name": fileName,
		},
	})
}

// LikeResource 点赞/取消点赞资源
func (c *ResourceController) LikeResource(ctx *gin.Context) {
	resourceID := ctx.Param("id")
	userID := ctx.GetString("user_id")

	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	var userLike models.UserLike
	result := c.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userLike)

	if result.Error == nil {
		// 取消点赞
		c.DB.Delete(&userLike)
		// 使用count查询获取当前点赞数
		var likeCount int64
		c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)
		ctx.JSON(http.StatusOK, gin.H{"isLiked": false, "likes": likeCount})
	} else {
		// 点赞
		userIDUint, _ := strconv.ParseUint(userID, 10, 32)
		resourceIDUint, _ := strconv.ParseUint(resourceID, 10, 32)
		userLike = models.UserLike{
			UserID:     uint(userIDUint),
			ResourceID: uint(resourceIDUint),
		}
		c.DB.Create(&userLike)
		// 使用count查询获取当前点赞数
		var likeCount int64
		c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)
		ctx.JSON(http.StatusOK, gin.H{"isLiked": true, "likes": likeCount})
	}
}

// FavoriteResource 收藏/取消收藏资源
func (c *ResourceController) FavoriteResource(ctx *gin.Context) {
	resourceID := ctx.Param("id")
	userID := ctx.GetString("user_id")

	var userFavorite models.UserFavorite
	result := c.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userFavorite)

	if result.Error == nil {
		// 取消收藏
		c.DB.Delete(&userFavorite)
		ctx.JSON(http.StatusOK, gin.H{"isFavorited": false})
	} else {
		// 收藏
		userIDUint, _ := strconv.ParseUint(userID, 10, 32)
		resourceIDUint, _ := strconv.ParseUint(resourceID, 10, 32)
		userFavorite = models.UserFavorite{
			UserID:     uint(userIDUint),
			ResourceID: uint(resourceIDUint),
		}
		c.DB.Create(&userFavorite)
		ctx.JSON(http.StatusOK, gin.H{"isFavorited": true})
	}
}

// SearchResources 搜索资源
func (c *ResourceController) SearchResources(ctx *gin.Context) {
	// 获取搜索参数
	query := ctx.Query("q")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	categoryID := ctx.Query("category")

	// 构建查询
	dbQuery := c.DB.Model(&models.Resource{}).
		Where("status = ?", "approved").
		Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")

	// 分类过滤
	if categoryID != "" {
		dbQuery = dbQuery.Where("category_id = ?", categoryID)
	}

	// 执行查询
	var resources []models.Resource
	var total int64

	dbQuery.Count(&total)
	dbQuery.Preload("User").Preload("Category").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resources)

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// CreateResource 创建资源
func (c *ResourceController) CreateResource(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title          string `json:"title" binding:"required"`
		Description    string `json:"description" binding:"required"`
		CategoryID     uint   `json:"category_id" binding:"required"`
		FilePath       string `json:"file_path" binding:"required"`
		FileSize       int64  `json:"file_size" binding:"required"`
		FileType       string `json:"file_type" binding:"required"`
		PointsRequired int    `json:"points_required"`
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

	// 创建资源记录
	resource := models.Resource{
		Title:          input.Title,
		Description:    input.Description,
		CategoryID:     input.CategoryID,
		FilePath:       input.FilePath,
		FileSize:       input.FileSize,
		FileType:       input.FileType,
		PointsRequired: input.PointsRequired,
		Status:         "pending", // 默认为待审核状态
		UserID:         userID.(uint),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if result := c.DB.Create(&resource); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建资源失败"})
		return
	}

	// 添加积分记录（上传资源奖励积分）
	pointRecord := models.PointRecord{
		UserID:      userID.(uint),
		Points:      10, // 上传资源奖励10积分
		Type:        "upload",
		ResourceID:  &resource.ID,
		Description: "上传资源奖励",
		CreatedAt:   time.Now(),
	}

	c.DB.Create(&pointRecord)

	// 更新用户积分
	c.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points + ?", 10))

	// 返回创建的资源
	ctx.JSON(http.StatusCreated, resource)
}

// UpdateResource 更新资源
func (c *ResourceController) UpdateResource(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取资源ID
	id := ctx.Param("id")

	// 查询资源
	var resource models.Resource
	result := c.DB.First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查是否是资源的拥有者
	if resource.UserID != userID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权修改此资源"})
		return
	}

	// 绑定请求数据
	var input struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		CategoryID     uint   `json:"category_id"`
		PointsRequired int    `json:"points_required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新资源
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Title != "" {
		updates["title"] = input.Title
	}

	if input.Description != "" {
		updates["description"] = input.Description
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

	if input.PointsRequired >= 0 {
		updates["points_required"] = input.PointsRequired
	}

	// 更新状态为待审核
	updates["status"] = "pending"

	// 保存更新
	if result := c.DB.Model(&resource).Updates(updates); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新资源失败"})
		return
	}

	// 重新查询资源以获取最新信息
	c.DB.Preload("User").Preload("Category").First(&resource, id)

	// 返回更新后的资源
	ctx.JSON(http.StatusOK, resource)
}

// DeleteResource 删除资源
func (c *ResourceController) DeleteResource(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取资源ID
	id := ctx.Param("id")

	// 查询资源
	var resource models.Resource
	result := c.DB.First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查是否是资源的拥有者或管理员
	var user models.User
	c.DB.First(&user, userID)

	if resource.UserID != userID.(uint) && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此资源"})
		return
	}

	// 删除资源文件
	ctxBg := context.Background()
	err := c.MinioClient.RemoveObject(ctxBg, config.GetEnv("MINIO_BUCKET", "pool"), resource.FilePath, minio.RemoveObjectOptions{})
	if err != nil {
		// 记录错误但继续删除数据库记录
		log.Printf("删除文件失败: %v", err)
	}

	// 删除资源记录
	if result := c.DB.Delete(&resource); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除资源失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "资源已删除"})
}

// GetUserResources 获取用户资源
func (c *ResourceController) GetUserResources(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	status := ctx.DefaultQuery("status", "")

	// 构建查询
	query := c.DB.Model(&models.Resource{}).Where("user_id = ?", userID)

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 执行查询
	var resources []models.Resource
	var total int64

	query.Count(&total)
	query.Preload("Category").
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resources)

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// UploadFile 上传文件
func (c *ResourceController) UploadFile(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件"})
		return
	}
	defer file.Close()

	// 生成唯一文件名
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := fmt.Sprintf("uploads/%d/%s", userID, fileName)

	// 上传到Minio
	ctxBg := context.Background()
	_, err = c.MinioClient.PutObject(ctxBg, config.GetEnv("MINIO_BUCKET", "pool"), filePath, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	// 返回文件信息
	ctx.JSON(http.StatusOK, gin.H{
		"file_path": filePath,
		"file_size": header.Size,
		"file_type": header.Header.Get("Content-Type"),
		"file_name": header.Filename,
	})
}

// DownloadFile 下载文件
func (c *ResourceController) DownloadFile(ctx *gin.Context) {
	// 获取资源ID
	id := ctx.Param("id")

	// 查询资源
	var resource models.Resource
	result := c.DB.First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查资源状态
	if resource.Status != "approved" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "资源未通过审核"})
		return
	}

	// 检查用户是否登录
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	// 检查用户积分是否足够
	var user models.User
	c.DB.First(&user, userID)

	// 如果不是资源拥有者且积分不足
	if resource.UserID != userID.(uint) && user.Points < resource.PointsRequired {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "积分不足"})
		return
	}

	// 获取文件
	ctxBg := context.Background()
	object, err := c.MinioClient.GetObject(ctxBg, config.GetEnv("MINIO_BUCKET", "pool"), resource.FilePath, minio.GetObjectOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件失败"})
		return
	}
	defer object.Close()

	// 获取文件信息
	stat, err := object.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件信息失败"})
		return
	}

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(resource.FilePath)))
	ctx.Header("Content-Type", stat.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", stat.Size))

	// 将文件内容写入响应
	ctx.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, nil)

	// 如果不是资源拥有者，扣除积分
	if resource.UserID != userID.(uint) && resource.PointsRequired > 0 {
		// 添加积分记录
		pointRecord := models.PointRecord{
			UserID:      userID.(uint),
			Points:      -resource.PointsRequired,
			Type:        "download",
			ResourceID:  &resource.ID,
			Description: "下载资源消费",
			CreatedAt:   time.Now(),
		}

		c.DB.Create(&pointRecord)

		// 更新用户积分
		c.DB.Model(&models.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points - ?", resource.PointsRequired))

		// 给资源上传者增加积分（分成机制）
		sharePoints := resource.PointsRequired / 2 // 上传者获得一半积分
		if sharePoints > 0 {
			// 添加积分记录
			ownerPointRecord := models.PointRecord{
				UserID:      resource.UserID,
				Points:      sharePoints,
				Type:        "share",
				ResourceID:  &resource.ID,
				Description: "资源被下载分成",
				CreatedAt:   time.Now(),
			}

			c.DB.Create(&ownerPointRecord)

			// 更新上传者积分
			c.DB.Model(&models.User{}).Where("id = ?", resource.UserID).Update("points", gorm.Expr("points + ?", sharePoints))
		}
	}

	// 更新下载次数
	c.DB.Model(&resource).Update("download_count", gorm.Expr("download_count + 1"))
}
