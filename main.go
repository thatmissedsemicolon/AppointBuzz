package main

import (
	"github.com/gin-gonic/gin"

	v1 "appointbuzz/api"
	db "appointbuzz/api/v1/lib"
)

func main() {
    router := gin.Default()

    db.InitializeDatabase()

    v1.SetupRouter(router)

    router.SetTrustedProxies(nil)

    router.Run(":8000")
}
