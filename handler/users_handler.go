package handler

import (
	"RestAPI-todo-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// Login handles user login and JWT token generation
func Login(c *gin.Context) {
	var data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if the user exists and the password matches
	user := models.UserMatchPassword(data.Username, data.Password)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// JWT token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"admin":    user.Admin,
		"ttl":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// Hard-Coded secret for demonstration purposes
	sampleSecret := []byte("secret")

	tokenString, err := token.SignedString(sampleSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return only the token string without setting a cookie
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
