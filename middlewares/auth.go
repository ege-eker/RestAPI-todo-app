package middlewares

import (
	"RestAPI-todo-app/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Validate username against mockUsers
		user := findUserByUsername(username.(string))
		if user == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func findUserByUsername(username string) *models.User {
	for _, user := range models.MockUsers() {
		if user.Username == username {
			return &user
		}
	}
	return nil
}
