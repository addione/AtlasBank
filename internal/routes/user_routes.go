package routes

import (
	"github.com/atlasbank/api/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(router *gin.RouterGroup, userController *controllers.UserController) {
	users := router.Group("/users")
	{
		users.GET("", userController.GetAllUsers)
		users.POST("", userController.CreateUser)
		users.GET("/:id", userController.GetUser)
	}
}
