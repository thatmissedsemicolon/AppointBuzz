package router

import (
	"appointbuzz/api/v1/services"
	"github.com/gin-gonic/gin"
)

func ConfigureAuthRoutes(group *gin.RouterGroup) {
	group.POST("/signup", services.CreateUserHandler)
	group.POST("/login", services.LoginUserHandler)
}

func ConfigureUserRoutes(group *gin.RouterGroup) {
	group.GET("/user", services.GetUserHandler)
	group.PATCH("/user", services.UpdateUserHandler)
	group.DELETE("/user", services.DeleteUserHandler)
}

func ConfigureAdminRoutes(group *gin.RouterGroup) {
	group.GET("/users", services.GetAllUsersHandler)
}
