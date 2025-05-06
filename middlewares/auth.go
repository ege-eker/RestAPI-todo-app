package middlewares

import (
	"RestAPI-todo-app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieving the JWT token from the cookie
		tokenStr, err := c.Cookie("Auth")
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Extracting the JWT token from the cookie
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Hard-Coded secret for demonstration purposes
			sampleSecret := []byte("secret")

			return sampleSecret, nil
		})
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Failed to parse token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "JWT claims are not valid"})
			c.Abort()
			return
		}
		// Check expiration time of the token
		if float64(time.Now().Unix()) > claims["ttl"].(float64) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Extracting the username from the token claims
		user := findUserByUsername(claims["username"].(string))
		if user == nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Setting the user in the context
		c.Set("user", user)

		// Admin set
		if claims["admin"].(bool) {
			c.Set("admin", true)
		} else {
			c.Set("admin", false)
		}

		c.Next()
	}
}

// Getting user by username
func findUserByUsername(username string) *models.User {
	for _, user := range models.MockUsers() {
		if user.Username == username {
			return &user
		}
	}
	return nil
}
