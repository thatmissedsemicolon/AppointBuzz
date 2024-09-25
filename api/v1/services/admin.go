package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "appointbuzz/api/v1/lib"
)

func GetAllUsers(c *gin.Context) {
    if _, ok, err := CheckUserPermissions(c, []string{"admin"}); !ok {
        c.IndentedJSON(http.StatusForbidden, gin.H{"error": err})
        return
    }

    var users []db.User
    if err := db.DB.Find(&users).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}
