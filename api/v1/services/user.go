package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "appointbuzz/api/v1/lib"
)

func GetUser(c *gin.Context) {
	email, ok, errMsg := CheckUserPermissions(c, []string{"user", "admin"})
	if !ok {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": errMsg})
		return
	}

	var user db.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
    email, ok, errMsg := CheckUserPermissions(c, []string{"user", "admin"})
    if !ok {
        c.IndentedJSON(http.StatusForbidden, gin.H{"error": errMsg})
        return
    }

    var user db.User
    if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

	if err := bindJSON(c, &user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    if err := db.DB.Save(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
