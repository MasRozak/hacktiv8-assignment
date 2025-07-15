package routes

import (
	"social-media-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", controllers.CreateUser)
		userRoutes.GET("", controllers.GetAllUsers)
		userRoutes.GET("/:id", controllers.GetUserByID)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.GET("/:id/posts", controllers.GetPostsByUserID)
		userRoutes.GET("/:id/likes", controllers.GetLikesByUserID)
		userRoutes.GET("/:id/followers", controllers.GetFollowers)
		userRoutes.GET("/:id/following", controllers.GetFollowing)
	}

	postRoutes := r.Group("/posts")
	{
		postRoutes.POST("", controllers.CreatePost)
		postRoutes.GET("", controllers.GetAllPosts)
		postRoutes.GET("/:id", controllers.GetPostByID)
		postRoutes.DELETE("/:id", controllers.DeletePost)
		postRoutes.GET("/:id/likes", controllers.GetLikesByPostID)
		postRoutes.GET("/:id/comments", controllers.GetCommentsByPostID)
	}

	likeRoutes := r.Group("/likes")
	{
		likeRoutes.POST("", controllers.CreateLike)
	}

	commentRoutes := r.Group("/comments")
	{
		commentRoutes.POST("", controllers.CreateComment)
	}

	followRoutes := r.Group("/follows")
	{
		followRoutes.POST("", controllers.CreateFollow)
		followRoutes.DELETE("", controllers.DeleteFollow)
	}
}
