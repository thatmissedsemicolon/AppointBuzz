package services

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	sql "appointbuzz/api/v1/sql"
)

func bindJSON(c *gin.Context, target interface{}) error {
    if err := c.BindJSON(target); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return err
    }
    return nil
}

func validateEmail(email string) (bool, error) {
    _, err := mail.ParseAddress(email)
    return err == nil, err
}

func userExists(email string) (bool, error) {
    var user sql.User
    err := sql.DB.Where("email = ?", email).First(&user).Error
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
