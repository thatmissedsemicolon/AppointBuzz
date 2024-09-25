package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	services "appointbuzz/api/v1/lib"
)

type UserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func CreateUserHandler(c *gin.Context) {
    var signupReq UserRequest
    if err := bindJSON(c, &signupReq); err != nil {
        return
    }

    _ , err := validateEmail(signupReq.Email)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
        return
    }

    exists, err := userExists(signupReq.Email)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }
    if exists {
        c.IndentedJSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    hashedPassword, err := hashPassword(signupReq.Password)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    newUser := services.User{Name: signupReq.Name, Email: signupReq.Email, Password: hashedPassword, Roles: convertRolesToString([]string{"user"})}
    if err := services.DB.Create(&newUser).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    accessToken, refreshToken, err := services.CreateTokens(newUser.Roles, newUser.Email)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.IndentedJSON(http.StatusCreated, gin.H{
        "access_token": accessToken,
        "refresh_token": refreshToken,
        "expires_in": 3600,
        "token_type": "Bearer",
    })
}

func LoginUserHandler(c *gin.Context) {
    var loginReq UserRequest
    if err := bindJSON(c, &loginReq); err != nil {
        return
    }

    exists, err := userExists(loginReq.Email)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }
    if !exists {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
        return
    }

    var user services.User
    if err := services.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
        return
    }

    accessToken, refreshToken, err := services.CreateTokens(user.Roles, user.Email)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{
        "access_token": accessToken,
        "refresh_token": refreshToken,
        "expires_in": 3600, 
        "token_type": "Bearer",
    })
}
