package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"appointbuzz/api/v1/lib"
)

func GetAllUsersHandler(c *gin.Context) {
	if _, ok, err := checkUserPermissions(c, []string{"admin"}); !ok {
		responseError(c, http.StatusForbidden, err)
		return
	}

	defaultLimit := 10
	defaultOffset := 0

	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(defaultOffset)))
	if err != nil {
		offset = defaultOffset
	}

	email := c.Query("email")

	var users []lib.User
	query := lib.DB.Offset(offset).Limit(limit)

	if email != "" {
		query = query.Where("email = ?", email)
	}

	if err := query.Find(&users).Error; err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}
