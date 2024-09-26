package services

import (
	"appointbuzz/api/v1/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func CreateUserHandler(c *gin.Context) {
    var signupReq UserRequest
    if err := c.ShouldBindJSON(&signupReq); err != nil {
        responseError(c, http.StatusBadRequest, "Invalid request data")
        return
    }

    if _, err := validateEmail(signupReq.Email); err != nil {
        responseError(c, http.StatusBadRequest, "Invalid email address")
        return
    }

    exists, isDeleted, err := userExists(signupReq.Email)
    if err != nil {
        responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    statusCode, message := handleExistingUserResponse(exists, isDeleted)
    if statusCode != 0 {
        responseError(c, statusCode, message)
        return
    }

    hashedPassword, _ := hashPassword(signupReq.Password)
    newUser := lib.User{Name: signupReq.Name, Email: signupReq.Email, Password: hashedPassword, Roles: convertRolesToString([]string{"user"})}
    if err := lib.DB.Create(&newUser).Error; err != nil {
        responseError(c, http.StatusInternalServerError, "Something went wrong")
        return
    }

    issueTokens(c, newUser)
}

func LoginUserHandler(c *gin.Context) {
    var loginReq UserRequest
    if err := c.ShouldBindJSON(&loginReq); err != nil {
        responseError(c, http.StatusBadRequest, "Invalid request data")
        return
    }

    var user lib.User
    if err := lib.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
        responseError(c, http.StatusUnauthorized, "Invalid login credentials")
        return
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
        responseError(c, http.StatusUnauthorized, "Incorrect email or password")
        return
    }

    issueTokens(c, user)
}
