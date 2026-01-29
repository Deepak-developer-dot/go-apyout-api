package controllers

import (
	"fmt"
	"net/http"
	"os"
	"payout-backend/config"
	"payout-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	Role   int  `json:"role"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println("Email from request", req.Email)
	// fmt.Println("Password from request", req.Password)

	// password := "1234568"
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Generated Hash:", string(hashedPassword))

	var merchant models.Merchant
	if err := config.DB.Where("email=?", req.Email).First(&merchant).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// fmt.Println("Input password:", req.Password)
	// fmt.Println("DB password hash:", merchant.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password1"})
		return
	}
	claims := &Claims{
		UserID: merchant.ID,
		Role:   int(merchant.IsAdmin),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	fmt.Println("Merchant Role:", merchant.IsAdmin)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could Not Login"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfull",
		"token":   tokenString,
		"user":    merchant,
	})
}

func Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	userID := c.GetUint("user_id")

	cacheKey := fmt.Sprintf("blacklist_token:%s", tokenString)

	config.Cache.Set(cacheKey, userID, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Logout Successfully!"})
}
