package main

import (
	"RestAPI-todo-app/handler"
	"RestAPI-todo-app/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	api.POST("/login", handler.Login)

	// API routes that require authentication middleware
	auth := api.Group("/", middlewares.AuthMiddleware())

	auth.GET("/todo-lists", handler.GetTodos)
	auth.POST("/todo-lists", handler.AddTodoList)
	auth.PUT("/todo-lists/:id", handler.RenameTodoList)
	auth.DELETE("/todo-lists/:id", handler.DeleteTodoList)
	auth.POST("/todo-lists/:id/step", handler.AddTodoStep)
	auth.DELETE("/todo-lists/:id/step/:step_id", handler.DeleteTodoStep)
	auth.PUT("/todo-lists/:id/step/:step_id", handler.RenameTodoStep)
	auth.POST("/todo-lists/:id/step/:step_id/toggle", handler.ToggleStepCompletion)

	// Gets Todo lists with deleted elements, only for admins. uncomment if needed
	//auth.GET("/admin/todolists", middlewares.AuthMiddleware(), handler.GetAllTodosAdmin)

	// Start the server on port 8080
	r.Run(":8080")
}
