package main

import (
	"github.com/gin-gonic/gin"

	"appointbuzz/api"
	"appointbuzz/api/v1/lib"
)

func main() {
	lib.InitRedis()

	lib.InitializeDatabase()

	router := gin.Default()
	v1.SetupRouter(router)
	router.SetTrustedProxies(nil)

	router.Run(":8000")
}
