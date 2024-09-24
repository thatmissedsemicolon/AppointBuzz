package main

import (
	"github.com/gin-gonic/gin"

	v1 "appointbuzz/api"
	sql "appointbuzz/api/v1/sql"
)

func main() {
    router := gin.Default()

    sql.InitializeDatabase()

    v1.SetupRouter(router)

    router.SetTrustedProxies(nil)

    router.Run(":8000")
}
