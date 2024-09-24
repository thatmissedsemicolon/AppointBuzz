package auth

import (
    "net/http"
    "net/mail"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    sql "appointbuzz/api/v1/sql"
    jwt "appointbuzz/lib"
)

type UserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

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
    err := sql.Db.Where("email = ?", email).First(&user).Error
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

    newUser := sql.User{ID: uuid.New(), Email: signupReq.Email, Password: hashedPassword}
    if err := sql.Db.Create(&newUser).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    token, err := jwt.CreateToken(newUser.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
        return
    }

    c.IndentedJSON(http.StatusCreated, gin.H{"access_token": token})
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
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User does not exists"})
        return
    }

    var user sql.User
    if err := sql.Db.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
        return
    }

    token, err := jwt.CreateToken(user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"access_token": token})
}

func GetAllUsers(c *gin.Context) {
    _, exists := c.Get("email")

    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
        return
    }

    var users []sql.User
    if err := sql.Db.Find(&users).Error; err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }
    c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}
