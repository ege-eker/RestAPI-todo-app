package controllers

import (
	"RestAPI-todo-app/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

type formData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(c *gin.Context) {
	var data formData
	if err := c.ShouldBind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid input"})
		return
	}

	user := models.UserMatchPassword(data.Username, data.Password)
	if user == nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()
	log.Println("session successfully saved")
	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	log.Println("session successfully cleared")
}
