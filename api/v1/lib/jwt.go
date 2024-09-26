package lib

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	loadKeysOnce   sync.Once
	privateKeyPath string = "private_key.pem"
	publicKeyPath  string = "public_key.pem"
)

type Claims struct {
	Email string `json:"email"`
	Roles string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func loadKeys() error {
	var err error
	loadKeysOnce.Do(func() {
		err = loadPrivateKey(privateKeyPath)
		if err != nil {
			err = fmt.Errorf("failed to load private key: %w", err)
			return
		}
		err = loadPublicKey(publicKeyPath)
		if err != nil {
			err = fmt.Errorf("failed to load public key: %w", err)
		}
	})
	return err
}

func loadPublicKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("public key block size is incorrect")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var ok bool
	if publicKey, ok = pub.(*rsa.PublicKey); !ok {
		return fmt.Errorf("not an RSA public key")
	}
	return nil
}

func loadPrivateKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("private key block size is incorrect")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	var ok bool
	if privateKey, ok = priv.(*rsa.PrivateKey); !ok {
		return fmt.Errorf("not an RSA private key")
	}
	return nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := loadKeys(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("email", claims.Email)
			c.Set("roles", claims.Roles)
			c.Next()
		}
	}
}

func CreateTokens(roles string, email string) (string, string, error) {
	if err := loadKeys(); err != nil {
		return "", "", err
	}

	accessTokenClaims := &Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", time.Now().Unix()),
			Issuer:    "https://appointbuzz.com",
			Audience:  jwt.ClaimStrings{"appointbuzz"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenClaims := &Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", time.Now().Unix()),
			Issuer:    "https://appointbuzz.com",
			Audience:  jwt.ClaimStrings{"appointbuzz"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}
