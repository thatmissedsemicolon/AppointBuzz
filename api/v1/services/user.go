package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"appointbuzz/api/v1/lib"
)

func GetUser(c *gin.Context) {
    email, ok, err := checkUserPermissions(c, []string{"user", "admin"})
    if !ok {
		responseError(c, http.StatusForbidden, err)
        return
    }

    var user lib.User
    if cachedUser, err := lib.GetValue(email); err == nil {
        if json.Unmarshal([]byte(cachedUser), &user) == nil {
            c.IndentedJSON(http.StatusOK, gin.H{"user": user})
            return
        }
    }

    if err := lib.DB.Where("email = ?", email).First(&user).Error; err != nil {
		responseError(c, http.StatusNotFound, "User not found")
        return
    }

    userJson, _ := json.Marshal(user)
    lib.SetValue(email, string(userJson), 10*time.Minute)
    c.IndentedJSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
    email, ok, err := checkUserPermissions(c, []string{"user", "admin"})
    if !ok {
		responseError(c, http.StatusForbidden, err)
        return
    }

    var user lib.User
    if err := lib.DB.Where("email = ?", email).First(&user).Error; err != nil {
		responseError(c, http.StatusNotFound, "User not found")
        return
    }

	if err := c.ShouldBindJSON(&user); err != nil {
        responseError(c, http.StatusBadRequest, "Invalid request data")
        return
    }

    if err := lib.DB.Save(&user).Error; err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
