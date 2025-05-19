package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"g/front/backend/config"
	"g/front/backend/controllers"
	"g/front/backend/middleware"
	"g/front/backend/migrations"
	"g/front/backend/routes"
	"g/front/backend/utils"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASSWORD", "123456"),
		config.GetEnv("DB_HOST", "47.121.210.209"),
		config.GetEnv("DB_PORT", "13306"),
		config.GetEnv("DB_NAME", "pool"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 运行数据库迁移
	migrations.RunMigrations(db)

	// 初始化基础数据
	utils.InitData(db)

	// 初始化Minio客户端
	minioClient, err := config.InitMinioClient()
	if err != nil {
		log.Fatalf("Minio客户端初始化失败: %v", err)
	}

	// 确保bucket存在
	config.EnsureBucketExists(minioClient, config.GetEnv("MINIO_BUCKET", "pool"))

	// 创建Gin实例
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 使用错误处理中间件
	r.Use(middleware.ErrorMiddleware())

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 设置404处理
	r.NoRoute(middleware.NotFoundHandler)

	// 创建Minio工具实例
	_ = utils.NewMinioUtils(minioClient)

	// 注册控制器
	userController := controllers.NewUserController(db, minioClient)
	resourceController := controllers.NewResourceController(db, minioClient)
	// 初始化Redis客户端
	redisClient, err := config.InitRedisClient()
	if err != nil {
		log.Fatalf("Redis客户端初始化失败: %v", err)
	}

	forumController := controllers.NewForumController(db, redisClient)
	chatController := controllers.NewChatController(db)
	pointsController := controllers.NewPointsController(db)
	adminController := controllers.NewAdminController(db)

	// 注册路由
	routes.SetupRoutes(r, userController, resourceController, forumController, chatController, pointsController, adminController)

	// 获取端口
	port := config.GetEnv("PORT", "8080")

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
