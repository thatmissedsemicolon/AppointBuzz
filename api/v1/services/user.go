package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sql "appointbuzz/api/v1/sql"
)

func GetAllUsers(c *gin.Context) {
    _, emailExists := c.Get("email"); roles, rolesExists := c.Get("roles"); 
    var rolesList []string; 
    if rolesExists { 
        rolesList = convertStringToRoles(roles.(string)) 
    }

    if contains(rolesList, "user") || contains(rolesList, "admin") {
        if !emailExists || !rolesExists {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
            return
        }
        
        var users []sql.User
        if err := sql.DB.Find(&users).Error; err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
            return
        }
        c.IndentedJSON(http.StatusOK, gin.H{"users": users})
    } else {
        c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You do not have the necessary permissions to access this resource."})
        return
    }
}
