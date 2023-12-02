package router

import (
	"upload-photos/controllers"
	"upload-photos/middlewares"

	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
	// Create a new router
	r := gin.Default()

	// Define the user routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", controllers.RegisterUser)
		userRoutes.POST("/login", controllers.LoginUser)
		userRoutes.PUT("/:userId", middlewares.AuthMiddleware, controllers.UpdateUser)
		userRoutes.DELETE("/:userId", middlewares.AuthMiddleware, controllers.DeleteUser)
	}

	// Define the photo routes
	photoRoutes := r.Group("/photos")
	{
		photoRoutes.POST("/", middlewares.AuthMiddleware, controllers.CreatePhoto)
		photoRoutes.GET("/", controllers.GetPhotos)
		photoRoutes.PUT("/:photoId", middlewares.AuthMiddleware, controllers.UpdatePhoto)
		photoRoutes.DELETE("/:photoId", middlewares.AuthMiddleware, controllers.DeletePhoto)
	}

	// Return the router
	return r
}
