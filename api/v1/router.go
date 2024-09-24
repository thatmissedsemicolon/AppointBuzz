package router

import (
    "github.com/gin-gonic/gin"
    auth "appointbuzz/api/v1/services"
)

func AuthRoutes(group *gin.RouterGroup) {
    group.POST("/signup", auth.CreateUserHandler)
    group.POST("/login", auth.LoginUserHandler)
}

func UserRoutes(group *gin.RouterGroup) {
    group.GET("/users", auth.GetAllUsers)
}
