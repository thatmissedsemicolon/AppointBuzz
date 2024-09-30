package router

import (
	"appointbuzz/api/v1/services"
	"github.com/gin-gonic/gin"
)

func ConfigureAuthRoutes(group *gin.RouterGroup) {
	group.POST("/signup", services.CreateUserHandler)
	group.POST("/login", services.LoginUserHandler)
}

func ConfigureFormRoutes(group *gin.RouterGroup) {
    group.POST("/:form_id", services.FormSubmissionHandler)
}

func ConfigureUserRoutes(group *gin.RouterGroup) {
	group.GET("/user", services.GetUserHandler)
	group.PATCH("/user", services.UpdateUserHandler)
	group.DELETE("/user", services.DeleteUserHandler)
    group.POST("/create-form", services.CreateFormHandler)
}

func ConfigureAdminRoutes(group *gin.RouterGroup) {
	group.GET("/users", services.GetAllUsersHandler)
}
