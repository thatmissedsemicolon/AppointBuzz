package router

import (
    "github.com/gin-gonic/gin"
    services "appointbuzz/api/v1/services"
)

func AuthRoutes(group *gin.RouterGroup) {
    group.POST("/signup", services.CreateUserHandler)
    group.POST("/login", services.LoginUserHandler)
}

func UserRoutes(group *gin.RouterGroup) {
    group.GET("/users", services.GetAllUsers)
}
