package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

// AddFavorite 添加资源收藏
func (c *ResourceController) AddFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取资源ID
	resourceID := ctx.Param("id")

	// 检查资源是否存在
	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查是否已收藏
	var existingFavorite models.UserFavorite
	if err := c.DB.Where("user_id = ? AND resource_id = ?", userID, resource.ID).First(&existingFavorite).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "已收藏该资源", "isFavorited": true})
		return
	}

	// 创建收藏记录
	favorite := models.UserFavorite{
		UserID:     userID.(uint),
		ResourceID: resource.ID,
		CreatedAt:  time.Now(),
	}

	if err := c.DB.Create(&favorite).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "isFavorited": true, "message": "收藏成功"})
}

// RemoveFavorite 取消资源收藏
func (c *ResourceController) RemoveFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取资源ID
	resourceID := ctx.Param("id")

	// 检查资源是否存在
	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 删除收藏记录
	if err := c.DB.Where("user_id = ? AND resource_id = ?", userID, resource.ID).Delete(&models.UserFavorite{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "已取消收藏", "isFavorited": false})
}

// GetFavoriteStatus 获取资源收藏状态
func (c *ResourceController) GetFavoriteStatus(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取资源ID
	resourceID := ctx.Param("id")

	// 检查资源是否存在
	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查是否已收藏
	var existingFavorite models.UserFavorite
	if err := c.DB.Where("user_id = ? AND resource_id = ?", userID, resource.ID).First(&existingFavorite).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"is_favorite": false, "isFavorited": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"is_favorite": true, "isFavorited": true})
}

// GetUserFavorites 获取用户收藏列表
func (c *ResourceController) GetUserFavorites(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 查询用户收藏的资源
	var favorites []models.UserFavorite
	if err := c.DB.Preload("Resource").Preload("Resource.User").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}

	// 格式化返回数据
	var resources []gin.H
	for _, favorite := range favorites {
		resources = append(resources, gin.H{
			"id":          favorite.Resource.ID,
			"title":       favorite.Resource.Title,
			"description": favorite.Resource.Description,
			"category":    favorite.Resource.Category,
			"user":        favorite.Resource.User,
			"created_at":  favorite.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resources})
}

// NewResourceController 创建资源控制器实例
func NewResourceController(db *gorm.DB, minioClient *minio.Client) *ResourceController {
	return &ResourceController{DB: db, MinioClient: minioClient}
}

// DeleteUserResource 删除用户资源
func (c *ResourceController) DeleteUserResource(ctx *gin.Context) {
	// 获取资源ID
	resourceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的资源ID"})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 检查资源是否存在且属于该用户
	var resource models.Resource
	if err := c.DB.Where("id = ? AND user_id = ?", resourceID, userID).First(&resource).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在或不属于当前用户"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询资源失败"})
		}
		return
	}

	// 删除资源
	if err := c.DB.Delete(&resource).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除资源失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// GetResources 获取资源列表
func (c *ResourceController) GetResources(ctx *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	categoryID := ctx.Query("category")
	sort := ctx.DefaultQuery("sort", "newest")
	// 获取搜索关键词
	query := ctx.Query("query")

	// 构建查询
	dbQuery := c.DB.Model(&models.Resource{}).Where("status = ?", "approved")

	// 关键词搜索
	if query != "" {
		query = strings.TrimSpace(query)
		if len(query) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词至少需要2个字符"})
			return
		}
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 分类过滤
	if categoryID != "" {
		dbQuery = dbQuery.Where("category_id = ?", categoryID)
	}

	// 排序
	switch sort {
	case "newest":
		dbQuery = dbQuery.Order("created_at DESC")
	case "popular":
		dbQuery = dbQuery.Order("download_count DESC")
	default:
		dbQuery = dbQuery.Order("created_at DESC")
	}

	// 执行查询
	var resources []models.Resource
	var total int64

	dbQuery.Count(&total)
	dbQuery.Preload("User").Preload("Category").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resources)

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"query":     ctx.Query("query"),                 // 返回原始查询关键词
		"category":  ctx.Query("category"),              // 返回分类ID
		"sort":      ctx.DefaultQuery("sort", "newest"), // 返回排序方式
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

// GetResourceDownloadUrl 获取资源下载URL
func (c *ResourceController) GetResourceDownloadUrl(ctx *gin.Context) {
	id := ctx.Param("id")

	var resource models.Resource
	result := c.DB.First(&resource, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查MinIO文件是否存在
	ctxBg := context.Background()
	_, err := c.MinioClient.StatObject(ctxBg, config.GetEnv("resources", "resources"), resource.FilePath, minio.StatObjectOptions{})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "文件不存在",
			"details": gin.H{
				"bucket": config.GetEnv("resources", "resources"),
				"key":    resource.FilePath,
			},
		})
		return
	}

	// 生成MinIO下载URL
	endpoint := config.GetEnv("MINIO_ENDPOINT", "47.121.210.209:9001")
	bucketName := config.GetEnv("resources", "resources")
	fileURL := fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, resource.FilePath)

	ctx.JSON(http.StatusOK, gin.H{
		"url":      fileURL,
		"filename": resource.FilePath,
	})
}

// GetCategories 获取资源分类
func (c *ResourceController) GetCategories(ctx *gin.Context) {
	var categories []models.Category
	c.DB.Find(&categories)

	ctx.JSON(http.StatusOK, categories)
}

// Search 搜索资源
// @Summary 搜索资源
// @Description 根据关键词、分类、标签等条件搜索资源
// @Tags resources
// @Accept json
// @Produce json
// @Param q query string false "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param tags query string false "标签，多个用逗号分隔"
// @Param sort query string false "排序方式" enums(created_at:desc,download_count:desc,title:asc,rating:desc) default(created_at:desc)
// @Param price_range query string false "价格范围" enums(all,free,paid) default(all)
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(12)
// @Success 200 {object} map[string]interface{} "{\"resources\": [], \"total\": 0}"
// @Router /resources/search [get]
func (c *ResourceController) Search(ctx *gin.Context) {
	// 获取查询参数
	q := ctx.Query("q")
	categoryID, _ := strconv.Atoi(ctx.Query("category_id"))
	tags := ctx.Query("tags")
	sort := ctx.Query("sort")
	priceRange := ctx.Query("price_range")
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	// 构建查询条件
	query := c.DB.Model(&models.Resource{}).Where("status = ?", "approved")

	// 关键词搜索
	if q != "" {
		q = strings.TrimSpace(q)
		if len(q) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词至少需要2个字符"})
			return
		}

		// 更精确的搜索条件
		query = query.Where(
			"MATCH(title, description) AGAINST(? IN BOOLEAN MODE) OR "+
				"title LIKE ? OR description LIKE ?",
			q+"*", "%"+q+"%", "%"+q+"%")
	}

	// 分类过滤
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 标签过滤
	if tags != "" {
		query = query.Joins("JOIN resource_tags ON resource_tags.resource_id = resources.id").
			Joins("JOIN tags ON tags.id = resource_tags.tag_id").
			Where("tags.name IN ?", strings.Split(tags, ","))
	}

	// 价格范围过滤
	if priceRange == "free" {
		query = query.Where("price = 0")
	} else if priceRange == "paid" {
		query = query.Where("price > 0")
	}

	// 排序
	switch sort {
	case "download_count:desc":
		query = query.Order("download_count DESC")
	case "title:asc":
		query = query.Order("title ASC")
	case "rating:desc":
		query = query.Order("rating DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 执行分页查询
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
		"query":     ctx.Query("query"),                 // 返回原始查询关键词
		"category":  ctx.Query("category"),              // 返回分类ID
		"sort":      ctx.DefaultQuery("sort", "newest"), // 返回排序方式
	})
}

// GetComments 获取资源评论
func (c *ResourceController) GetComments(ctx *gin.Context) {
	// 获取资源ID
	resourceID := ctx.Param("id")
	if resourceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "资源ID不能为空"})
		return
	}

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

func (c *ResourceController) DeleteComment(ctx *gin.Context) {
	// 获取评论ID
	commentID := ctx.Param("id")

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询评论
	var comment models.Comment
	if err := c.DB.First(&comment, commentID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限：评论所有者或管理员
	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if comment.UserID != userID.(uint) && user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此评论"})
		return
	}

	// 删除评论
	if err := c.DB.Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}

// PostComment 提交资源评论
func (c *ResourceController) PostComment(ctx *gin.Context) {
	// 获取资源ID
	resourceID := ctx.Param("id")
	if resourceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "资源ID不能为空"})
		return
	}

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
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

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
		resourceIDUint, _ := strconv.ParseUint(resourceID, 10, 32)
		userLike = models.UserLike{
			UserID:     userID.(uint),
			ResourceID: uint(resourceIDUint),
		}
		c.DB.Create(&userLike)
		// 使用count查询获取当前点赞数
		var likeCount int64
		c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)
		ctx.JSON(http.StatusOK, gin.H{"isLiked": true, "likes": likeCount})
	}
}

// DislikeResource 取消资源点赞
func (c *ResourceController) DislikeResource(ctx *gin.Context) {
	resourceID := ctx.Param("id")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	var userLike models.UserLike
	result := c.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userLike)

	if result.Error == nil {
		// 存在点赞记录，删除它
		c.DB.Delete(&userLike)
		// 使用count查询获取当前点赞数
		var likeCount int64
		c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)
		ctx.JSON(http.StatusOK, gin.H{"isLiked": false, "likes": likeCount})
	} else {
		// 没有点赞记录，返回当前状态
		var likeCount int64
		c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)
		ctx.JSON(http.StatusOK, gin.H{"isLiked": false, "likes": likeCount})
	}
}

// GetResourceLikeStatus 获取资源点赞状态
func (c *ResourceController) GetResourceLikeStatus(ctx *gin.Context) {
	resourceID := ctx.Param("id")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var resource models.Resource
	if err := c.DB.First(&resource, resourceID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 检查用户是否已点赞
	var userLike models.UserLike
	result := c.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userLike)
	isLiked := result.Error == nil

	// 获取点赞总数
	var likeCount int64
	c.DB.Model(&models.UserLike{}).Where("resource_id = ?", resourceID).Count(&likeCount)

	ctx.JSON(http.StatusOK, gin.H{"isLiked": isLiked, "likes": likeCount})
}

// FavoriteResource 收藏/取消收藏资源
func (c *ResourceController) FavoriteResource(ctx *gin.Context) {
	resourceID := ctx.Param("id")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var userFavorite models.UserFavorite
	result := c.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userFavorite)

	if result.Error == nil {
		// 取消收藏
		c.DB.Delete(&userFavorite)
		ctx.JSON(http.StatusOK, gin.H{"isFavorited": false})
	} else {
		// 收藏
		userIDStr, ok := userID.(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
			return
		}
		userIDUint, _ := strconv.ParseUint(userIDStr, 10, 32)
		resourceIDUint, _ := strconv.ParseUint(resourceID, 10, 32)
		userFavorite = models.UserFavorite{
			UserID:     uint(userIDUint),
			ResourceID: uint(resourceIDUint),
		}
		c.DB.Create(&userFavorite)
		ctx.JSON(http.StatusOK, gin.H{"isFavorited": true})
	}
}

// DeleteMyResource 删除用户资源
func (c *ResourceController) DeleteMyResource(ctx *gin.Context) {
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此资源"})
		return
	}

	// 删除资源
	if result := c.DB.Delete(&resource); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除资源失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// GetMyResources 获取当前用户的资源列表
func (c *ResourceController) GetMyResources(ctx *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 构建查询
	dbQuery := c.DB.Model(&models.Resource{}).
		Where("user_id = ?", userID)

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

// SearchResources 搜索资源
func (c *ResourceController) SearchResources(ctx *gin.Context) {
	// 获取搜索参数
	query := ctx.Query("query")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	categoryID := ctx.Query("category")

	// 构建查询
	dbQuery := c.DB.Model(&models.Resource{}).
		Where("status = ?", "approved")

	// 关键词搜索
	if query != "" {
		query = strings.TrimSpace(query)
		if len(query) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词至少需要2个字符"})
			return
		}
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

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
		"query":     ctx.Query("query"),                 // 返回原始查询关键词
		"category":  ctx.Query("category"),              // 返回分类ID
		"sort":      ctx.DefaultQuery("sort", "newest"), // 返回排序方式
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
	err := c.MinioClient.RemoveObject(ctxBg, config.GetEnv("MINIO_BUCKET", "resources"), resource.FilePath, minio.RemoveObjectOptions{})
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
		"query":     ctx.Query("query"),                 // 返回原始查询关键词
		"category":  ctx.Query("category"),              // 返回分类ID
		"sort":      ctx.DefaultQuery("sort", "newest"), // 返回排序方式
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

	// 验证用户权限
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

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

	// 检查资源所有者
	if resource.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权下载该资源"})
		return
	}

	// 增加下载次数
	c.DB.Model(&resource).Update("download_count", gorm.Expr("download_count +?", 1))

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
		log.Printf("获取文件信息失败: %v, 文件路径: %s, 桶名: %s", err, resource.FilePath, config.GetEnv("MINIO_BUCKET", "pool"))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件信息失败", "details": err.Error()})
		return
	}

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(resource.FilePath)))
	ctx.Header("Content-Type", stat.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", stat.Size))

	// 将文件内容写入响应
	ctx.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, nil)

	// 更新下载次数
	c.DB.Model(&resource).Update("download_count", gorm.Expr("download_count + 1"))
}
