package services

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

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

func formatDuration(d time.Duration) string {
    min := d / time.Minute
    sec := (d % time.Minute) / time.Second
    if min > 0 {
        return fmt.Sprintf("%d min %d sec", min, sec)
    }
    return fmt.Sprintf("%d sec", sec)
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
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600,
		"token_type":    "Bearer",
	})
}

func parseFormData(c *gin.Context) (map[string]interface{}, error) {
	contentType := c.GetHeader("Content-Type")
	var formData map[string]interface{}

	switch {
	case contentType == "application/json":
		if err := c.ShouldBindJSON(&formData); err != nil {
			return nil, err
		}
	case contentType == "application/x-www-form-urlencoded", contentType == "multipart/form-data":
		if err := c.Request.ParseForm(); err != nil {
			return nil, err
		}
		formData = make(map[string]interface{})
		for key, values := range c.Request.PostForm {
			formData[key] = values[0]
		}
	default:
		return nil, fmt.Errorf("unsupported content type")
	}

	if err := validateAndModifyFormData(formData); err != nil {
		return nil, err
	}

	return formData, nil
}

func validateAndModifyFormData(formData map[string]interface{}) error {
	if email, ok := formData["email"].(string); ok {
		if valid, err := validateEmail(email); !valid || err != nil {
			return errors.New("invalid email format")
		}
	}

	if name, ok := formData["name"].(string); ok {
		if len(name) < 2 || len(name) > 100 {
			return errors.New("name must be between 2 and 100 characters")
		}
	}

	if password, ok := formData["password"].(string); ok {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return err
		}
		formData["password"] = string(hashedPassword)
	}

	if message, ok := formData["message"].(string); ok {
		sanitizedMessage, err := sanitizeText(message)
		if err != nil {
			return err
		}
		formData["message"] = sanitizedMessage
	}

	return nil
}

func sanitizeText(input string) (string, error) {
    badWords := []string{"badword1", "badword2", "badword3"}
    lowerInput := strings.ToLower(input)

    for _, word := range badWords {
        if strings.Contains(lowerInput, word) {
            cleanWord := strings.Repeat("*", len(word))
            lowerInput = strings.ReplaceAll(lowerInput, word, cleanWord)
        }
    }

    return lowerInput, nil
}

func responseError(c *gin.Context, status int, message string) {
	c.IndentedJSON(status, gin.H{"error": message})
}
