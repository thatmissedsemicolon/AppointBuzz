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
        router.ConfigureAuthRoutes(authGroup)
    }

    protectedRoutes := api.Group("/")
    protectedRoutes.Use(lib.JWTAuthMiddleware())
    {
        router.ConfigureUserRoutes(protectedRoutes)
    }

    adminProtectedRoutes := api.Group("/admin")
    adminProtectedRoutes.Use(lib.JWTAuthMiddleware())
    {
        router.ConfigureAdminRoutes(adminProtectedRoutes)
    }

    return api
}
