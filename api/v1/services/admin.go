package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"appointbuzz/api/v1/lib"
)

func GetAllUsers(c *gin.Context) {
    if _, ok, err := checkUserPermissions(c, []string{"admin"}); !ok {
		responseError(c, http.StatusForbidden, err)
        return
    }

    var users []lib.User
    if err := lib.DB.Find(&users).Error; err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}
