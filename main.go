package main

import (
    "context"

    "github.com/gin-gonic/gin"

    db "appointbuzz/api/v1/lib"
    redis "appointbuzz/api/v1/lib"
    v1 "appointbuzz/api"
)

var ctx = context.Background()

func main() {
    redis.InitRedis()

    db.InitializeDatabase()

    router := gin.Default()
    v1.SetupRouter(router)
    router.SetTrustedProxies(nil)

    router.Run(":8000")
}
