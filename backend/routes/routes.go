package routes

import (
	"github.com/gin-gonic/gin"

	"g/front/backend/controllers"
	"g/front/backend/middleware"
	"g/front/backend/models"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine, userController *controllers.UserController, resourceController *controllers.ResourceController, forumController *controllers.ForumController, chatController *controllers.ChatController, pointsController *controllers.PointsController, adminController *controllers.AdminController) {
	// API路由组
	api := r.Group("/api")

	// 公开路由
	public := api.Group("/")
	{
		// 用户认证
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)

		// 资源相关路由
		resourceRoutes := public.Group("/resources")
		{
			resourceRoutes.GET("", resourceController.GetResources)
			resourceRoutes.GET("/categories", resourceController.GetCategories)
			resourceRoutes.GET("/:id", resourceController.GetResourceById)
			resourceRoutes.GET("/search", resourceController.SearchResources)

		}

		// 资源评论

		// 论坛相关
		public.GET("/forum/categories", forumController.GetCategories)
		public.GET("/forum/topics", forumController.GetTopics)
		public.GET("/forum/topics/:id", forumController.GetTopicById)
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/user/profile", userController.GetProfile)
		protected.PUT("/user/profile", userController.UpdateProfile)

		// 积分相关
		protected.GET("/user/points", pointsController.GetUserPoints)
		protected.GET("/user/points/history", pointsController.GetPointsHistory)

		// 资源管理
		protected.POST("/resources/:id/comments", resourceController.PostComment)
		protected.POST("/resources", resourceController.CreateResource)
		protected.PUT("/resources/:id", resourceController.UpdateResource)
		protected.DELETE("/resources/:id", resourceController.DeleteResource)
		protected.GET("/user/resources", resourceController.GetUserResources)
		protected.POST("/resources/upload", resourceController.UploadResource)
		protected.GET("/download/:id", resourceController.DownloadFile)

		// 论坛管理
		protected.POST("/forum/topics", forumController.CreateTopic)
		protected.PUT("/forum/topics/:id", forumController.UpdateTopic)
		protected.DELETE("/forum/topics/:id", forumController.DeleteTopic)
		protected.POST("/forum/topics/:id/replies", forumController.CreateReply)
		protected.PUT("/forum/replies/:id", forumController.UpdateReply)
		protected.DELETE("/forum/replies/:id", forumController.DeleteReply)

		// 管理员路由
		admin := protected.Group("/admin")
		// 管理员权限检查中间件
		admin.Use(func(c *gin.Context) {
			// 获取用户ID
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(401, gin.H{"error": "未授权"})
				c.Abort()
				return
			}

			// 查询用户角色
			var user models.User
			result := adminController.DB.First(&user, userID)
			if result.Error != nil {
				c.JSON(404, gin.H{"error": "用户不存在"})
				c.Abort()
				return
			}

			// 检查是否为管理员
			if user.Role != "admin" {
				c.JSON(403, gin.H{"error": "权限不足"})
				c.Abort()
				return
			}

			c.Next()
		})
		{
			// 资源审核
			admin.GET("/resources/pending", adminController.GetPendingResources)
			admin.PUT("/resources/:id/review", adminController.ReviewResource)

			// 积分管理
			admin.POST("/points/add", pointsController.AddPoints)

			// 统计信息
			admin.GET("/stats", adminController.GetUserStats)
		}

		// AI聊天
		protected.GET("/chat/histories", chatController.GetChatHistories)
		protected.GET("/chat/histories/:id", chatController.GetChatHistoryById)
		protected.POST("/chat/histories", chatController.CreateChatHistory)
		protected.PUT("/chat/histories/:id", chatController.UpdateChatHistory)
		protected.DELETE("/chat/histories/:id", chatController.DeleteChatHistory)
		protected.POST("/chat/histories/:id/messages", chatController.SendMessage)
	}
}
