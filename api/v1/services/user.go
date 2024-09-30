package services

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
		responseError(c, http.StatusInternalServerError, "Update failed: "+result.Error.Error())
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

func CreateFormHandler(c *gin.Context) {
	email, ok, err := checkUserPermissions(c, []string{"user"})
	if !ok {
		responseError(c, http.StatusForbidden, err)
		return
	}

	var formConfig lib.FormConfig
	if err := c.ShouldBindJSON(&formConfig); err != nil {
		responseError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	formConfig.ID = uuid.New()
	formConfig.UserEmail = email
	if err := lib.DB.Create(&formConfig).Error; err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Form configuration created successfully", "form_id": formConfig.ID})
}

func FormSubmissionHandler(c *gin.Context) {
    formConfigID := c.Param("form_id")
    parsedFormConfigID, err := uuid.Parse(formConfigID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form configuration ID"})
        return
    }
    domainName := c.Request.Host
    
    var formConfig lib.FormConfig
    if err := lib.DB.Where("id = ?", parsedFormConfigID).First(&formConfig).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Form configuration not found"})
        return
    }

    domains := strings.Split(formConfig.AllowedDomains, ",")
    if !contains(domains, domainName) {
        c.JSON(http.StatusForbidden, gin.H{"error": "This domain is not allowed for this form"})
        return
    }

    formData, err := parseFormData(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    jsonData, err := json.Marshal(formData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing form data"})
        return
    }

    newSubmissionID := uuid.New()
    form := lib.Form{
        ID:           newSubmissionID,
        FormConfigID: parsedFormConfigID,
        Data:         string(jsonData),
    }

    if err := lib.DB.Create(&form).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save form data"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Form submitted successfully"})
}
