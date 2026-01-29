package middleware

import (
	"fmt"
	"net/http"
	"os"
	"payout-backend/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SCERET"))

type Claims struct {
	UserID uint `json:"user_id"`
	Role   int  `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header Missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Missing"})
			c.Abort()
			return
		}

		if _, found := config.Cache.Get("blacklist_token:" + tokenString); found {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Expired, please login"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if !token.Valid || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalied or Expired Token"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*Claims); ok {
			fmt.Println("Role Set:", claims.Role)
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}
