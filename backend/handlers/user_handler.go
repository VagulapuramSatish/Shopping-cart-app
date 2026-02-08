package handlers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key") // change this to a secure secret in production

// CreateUser handles POST /users for user signup
func CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if user exists
	var existing models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashed),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID":        user.ID,
		"Username":  user.Username,
		"CartID":    user.CartID,
		"CreatedAt": user.CreatedAt,
	})
}

// Login handles POST /users/login for user authentication
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	expiration := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save token to user (single-session token)
	user.Token = signedToken
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"ID":       user.ID,
		"Username": user.Username,
		"Token":    signedToken,
	})
}

// ListUsers handles GET /users to list all users
func ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// hide passwords in the response
	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"ID":        u.ID,
			"Username":  u.Username,
			"CartID":    u.CartID,
			"CreatedAt": u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}
