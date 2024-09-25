package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "appointbuzz/api/v1/lib"
)

func GetAllUsers(c *gin.Context) {
    _, emailExists := c.Get("email")
    roles, rolesExists := c.Get("roles")
    
    if !emailExists || !rolesExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user information"})
        return
    }

    rolesList := convertStringToRoles(roles.(string))
    if !contains(rolesList, "user") && !contains(rolesList, "admin") {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        return
    }

    var users []db.User
    if err := db.DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}
