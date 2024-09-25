package router

import (
    "github.com/gin-gonic/gin"
    "appointbuzz/api/v1/services"
)

func AuthRoutes(group *gin.RouterGroup) {
    group.POST("/signup", services.CreateUserHandler)
    group.POST("/login", services.LoginUserHandler)
}

func UserRoutes(group *gin.RouterGroup) {
    group.GET("/user", services.GetUser)
    group.PATCH("/user/update", services.UpdateUser)
}

func AdminRoutes(group *gin.RouterGroup) {
    group.GET("/users", services.GetAllUsers)
}
