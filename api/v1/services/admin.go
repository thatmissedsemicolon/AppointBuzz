package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "appointbuzz/api/v1/lib"
)

func GetAllUsers(c *gin.Context) {
    if _, ok, err := checkUserPermissions(c, []string{"admin"}); !ok {
		responseError(c, http.StatusForbidden, err)
        return
    }

    var users []db.User
    if err := db.DB.Find(&users).Error; err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}
