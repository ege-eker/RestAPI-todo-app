package main

import (
	"RestAPI-todo-app/controllers"
	"RestAPI-todo-app/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	r.GET("/", middlewares.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run(":8080")
}
