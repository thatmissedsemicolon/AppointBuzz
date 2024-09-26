package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"appointbuzz/api/v1/lib"
)

func GetUserHandler(c *gin.Context) {
    email, ok, err := checkUserPermissions(c, []string{"user"})
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

func UpdateUserHandler(c *gin.Context) {
    email, ok, err := checkUserPermissions(c, []string{"user"})
    if !ok {
		responseError(c, http.StatusForbidden, err)
        return
    }

    var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

    if password, exists := updates["password"]; exists {
		hashedPassword, err := hashPassword(password.(string))
		if err != nil {
			responseError(c, http.StatusInternalServerError, "Something went wrong")
			return
		}
		updates["password"] = string(hashedPassword)
	}

	result := lib.DB.Model(&lib.User{}).Where("email = ?", email).Updates(updates)
    if result.Error != nil {
        responseError(c, http.StatusInternalServerError, "Update failed: " + result.Error.Error())
        return
    }

    if result.RowsAffected == 0 {
        responseError(c, http.StatusNotFound, "No user found or no new data provided")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUserHandler(c *gin.Context) {
    email, ok, err := checkUserPermissions(c, []string{"user"})
    if !ok {
        responseError(c, http.StatusForbidden, err)
        return
    }

    result := lib.DB.Where("email = ?", email).Delete(&lib.User{})
    if result.Error != nil {
        responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    if result.RowsAffected == 0 {
        responseError(c, http.StatusNotFound, "No user found")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
