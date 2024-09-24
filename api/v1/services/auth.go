package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	sql "appointbuzz/api/v1/sql"
	jwt "appointbuzz/lib"
)

type UserRequest struct {
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
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if exists {
        c.IndentedJSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    hashedPassword, err := hashPassword(signupReq.Password)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    newUser := sql.User{Email: signupReq.Email, Password: hashedPassword, Roles: convertRolesToString([]string{"user"})}
    if err := sql.DB.Create(&newUser).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    accessToken, refreshToken, err := jwt.CreateTokens(newUser.Roles, newUser.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
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
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if !exists {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
        return
    }

    var user sql.User
    if err := sql.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
        return
    }

    accessToken, refreshToken, err := jwt.CreateTokens(user.Roles, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{
        "access_token": accessToken,
        "refresh_token": refreshToken,
        "expires_in": 3600, 
        "token_type": "Bearer",
    })
}
