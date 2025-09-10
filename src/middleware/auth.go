package middleware

import (
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey string

func SetSecretKey(key string) {
	secretKey = key
}

func GetSecretKey() string {
	return secretKey
}

func AuthMiddleware() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		// Gets the authorization header
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		// Divides the header into Bearer and Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			ctx.Abort()
			return
		}

		// Verifies the JWT token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		// Checks if the token is valid
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Adds expiration validation for the token
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				ctx.Abort()
				return
			}
		}

		// Sets the token claims in the context (user ID)
		ctx.Set("userId", claims["id"])
		ctx.Next()
	}
}
