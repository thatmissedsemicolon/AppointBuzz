package services

import (
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	db "appointbuzz/api/v1/lib"
)

func validateEmail(email string) (bool, error) {
    _, err := mail.ParseAddress(email)
    return err == nil, err
}

func userExists(email string) (bool, error) {
    var user db.User
    err := db.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(bytes), err
}

func convertRolesToString(roles []string) string {
    return strings.Join(roles, ",")
}

func convertStringToRoles(roles string) []string {
    if roles == "" {
        return []string{}
    }
    return strings.Split(roles, ",")
}

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func checkUserPermissions(c *gin.Context, requiredRoles []string) (string, bool, string) {
    email, emailExists := c.Get("email")
    roles, rolesExists := c.Get("roles")

    if !emailExists || !rolesExists {
        return email.(string), false, "Missing user information"
    }

    rolesList := convertStringToRoles(roles.(string))
    for _, role := range requiredRoles {
        if contains(rolesList, role) {
            return email.(string), true, ""
        }
    }

    return email.(string), false, "Insufficient permissions"
}

func responseError(c *gin.Context, status int, message string) {
	c.IndentedJSON(status, gin.H{"error": message})
}
