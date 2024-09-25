package v1

import (
    "github.com/gin-gonic/gin"
    "appointbuzz/api/v1"
    "appointbuzz/api/v1/lib"
)

func SetupRouter(route *gin.Engine) *gin.RouterGroup {
    api := route.Group("/api/v1")

    authGroup := api.Group("/auth")
    {
        router.AuthRoutes(authGroup)
    }

    protectedRoutes := api.Group("/")
    protectedRoutes.Use(lib.JWTAuthMiddleware())
    {
        router.UserRoutes(protectedRoutes)
    }

    adminProtectedRoutes := api.Group("/admin")
    adminProtectedRoutes.Use(lib.JWTAuthMiddleware())
    {
        router.AdminRoutes(adminProtectedRoutes)
    }

    return api
}
