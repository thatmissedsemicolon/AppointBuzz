package services

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"appointbuzz/api/v1/lib"
)

func validateEmail(email string) (bool, error) {
    _, err := mail.ParseAddress(email)
    return err == nil, err
}

func userExists(email string) (exists bool, isDeleted bool, err error) {
    var user lib.User
    err = lib.DB.Unscoped().Where("email = ?", email).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return false, false, nil
        }
        return false, false, err
    }
    return true, !user.DeletedAt.Time.IsZero(), nil
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

func handleExistingUserResponse(exists bool, isDeleted bool) (int, string) {
    if exists {
        if isDeleted {
            return http.StatusGone, "This account has been permanently deleted. If you believe this was a mistake or require further assistance, please contact support."
        }
        return http.StatusConflict, "User already exists."
    }
    return 0, ""
}

func issueTokens(c *gin.Context, user lib.User) {
    accessToken, refreshToken, err := lib.CreateTokens(user.Roles, user.Email)
    if err != nil {
        responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{
        "access_token": accessToken,
        "refresh_token": refreshToken,
        "expires_in": 3600,
        "token_type": "Bearer",
    })
}

func responseError(c *gin.Context, status int, message string) {
	c.IndentedJSON(status, gin.H{"error": message})
}
