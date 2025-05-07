package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		tokenStr := ""

		if authHeader != "" {
			// Check if it has the "Bearer " prefix
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenStr = authHeader[7:]
			}
		}

		// If no token found in header, return unauthorized
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No valid token found"})
			c.Abort()
			return
		}

		// Validate the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Hard-Coded secret for demonstration purposes
			sampleSecret := []byte("secret")

			return sampleSecret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT claims are not valid"})
			c.Abort()
			return
		}
		// Check expiration time of the token
		if float64(time.Now().Unix()) > claims["ttl"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Setting the username in the context
		c.Set("username", claims["username"].(string))

		// Admin set
		c.Set("admin", claims["admin"].(bool))

		c.Next()
	}
}
