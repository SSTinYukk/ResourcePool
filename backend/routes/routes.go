package routes

import (
	"github.com/gin-gonic/gin"

	"g/front/backend/controllers"
	"g/front/backend/middleware"
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
			resourceRoutes.GET("/:id/comments", resourceController.GetComments)
		}

		// 资源评论

		// 论坛相关
		public.GET("/forum/categories", forumController.GetCategories)
		public.GET("/forum/topics", forumController.GetTopics)
		public.GET("/forum/topics/:id", forumController.GetTopicById)
		public.GET("/forum/topics/:id/likes", forumController.GetTopicLikes)
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/user/profile", userController.GetUserProfile)
		protected.PUT("/user/profile", userController.UpdateProfile)
		protected.PUT("/user/avatar", userController.UploadAvatar)
		protected.POST("/user/avatar", userController.UploadAvatar)
		protected.POST("/user/refresh-token", userController.RefreshToken)
		// 评论管理
		protected.DELETE("/resources/comments/:id", resourceController.DeleteComment)

		// 用户资源管理
		protected.DELETE("/user/resources/:id", resourceController.DeleteUserResource)
		protected.GET("/user/favorites/posts", forumController.GetUserFavorites)
		protected.GET("/user/favorites/resources", resourceController.GetUserFavorites)

		// 积分相关
		protected.GET("/user/points", pointsController.GetUserPoints)
		protected.GET("/user/points/history", pointsController.GetPointsHistory)

		// 资源管理
		protected.POST("/resources/:id/comments", resourceController.PostComment)
		protected.POST("/resources", resourceController.CreateResource)
		protected.PUT("/resources/:id", resourceController.UpdateResource)
		protected.DELETE("/resources/:id", resourceController.DeleteResource)
		protected.GET("/user/resources", resourceController.GetUserResources)
		protected.GET("/user/my-resources", resourceController.GetMyResources)
		protected.DELETE("/user/my-resources/:id", resourceController.DeleteMyResource)
		// 资源收藏
		protected.POST("/resources/:id/favorite", resourceController.AddFavorite)
		protected.DELETE("/resources/:id/favorite", resourceController.RemoveFavorite)

		protected.GET("/resources/:id/favorite-status", resourceController.GetFavoriteStatus)
		protected.GET("/favorites/resources", resourceController.GetUserFavorites)
		protected.POST("/resources/upload", resourceController.UploadResource)
		protected.GET("/download/:id", resourceController.GetResourceDownloadUrl)

		// 资源点赞
		protected.POST("/resources/:id/like", resourceController.LikeResource)
		protected.DELETE("/resources/:id/dislike", resourceController.DislikeResource)
		protected.GET("/resources/:id/like-status", resourceController.GetResourceLikeStatus)

		// 论坛管理
		protected.POST("/forum/topics", forumController.CreateTopic)
		protected.PUT("/forum/topics/:id", forumController.UpdateTopic)
		protected.DELETE("/forum/topics/:id", forumController.DeleteTopic)
		protected.POST("/forum/topics/:id/replies", forumController.CreateReply)
		protected.PUT("/forum/replies/:id", forumController.UpdateReply)
		protected.DELETE("/forum/replies/:id", forumController.DeleteReply)
		protected.POST("/forum/topics/:id/like", forumController.LikeTopic)
		protected.DELETE("/forum/topics/:id/like", forumController.UnlikeTopic)
		protected.POST("/forum/topics/:id/dislike", forumController.DislikeTopic)
		protected.DELETE("/forum/topics/:id/dislike", forumController.UnDislikeTopic)

		// 论坛收藏
		protected.POST("/forum/topics/:id/favorite", forumController.AddFavorite)
		protected.DELETE("/forum/topics/:id/favorite", forumController.RemoveFavorite)
		protected.GET("/forum/topics/:id/favorite-status", forumController.GetFavoriteStatus)

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			// 资源审核
			admin.GET("/resources/pending", adminController.GetPendingResources)
			admin.PUT("/resources/:id/review", adminController.ReviewResource)

			// 用户管理
			admin.GET("/users", adminController.GetUsers)
			admin.DELETE("/users/:id", adminController.DeleteUser)
			admin.PUT("/users/:id/role", adminController.UpdateUserRole)

			// 资源管理
			admin.GET("/resources", adminController.GetResources)
			admin.DELETE("/resources/:id", adminController.DeleteResource)

			// 论坛管理
			admin.GET("/forum/topics", adminController.GetTopics)
			admin.DELETE("/forum/topics/:id", adminController.DeleteTopic)

			// 积分管理
			admin.GET("/points/records", adminController.GetPointsRecords)
			admin.POST("/points/adjust", adminController.AdjustPoints)

			// 统计信息
			admin.GET("/stats", adminController.GetUserStats)
			admin.GET("/stats/users", adminController.GetUserStats)
			admin.GET("/stats/resources", adminController.GetResourceStats)
			admin.GET("/stats/forum", adminController.GetForumStats)
		}

		// AI聊天
		chat := protected.Group("/chat")
		{
			chat.POST("/sessions", chatController.CreateSession)
			chat.GET("/sessions", chatController.GetSessions)
			chat.GET("/sessions/:id/messages", chatController.GetSessionMessages)
			chat.DELETE("/sessions/:id", chatController.DeleteSession)
			chat.POST("/messages", chatController.SendMessage)
		}
	}
}
