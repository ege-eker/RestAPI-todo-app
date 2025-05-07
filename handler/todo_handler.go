package handler

import (
	"RestAPI-todo-app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// getUsername retrieves the username from the context
func getUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// getAdminState retrieves the admin state from the context
func getAdminState(c *gin.Context) bool {
	role, exists := c.Get("admin")
	if !exists {
		return false
	}
	return role.(bool)
}

// GetTodos retrieves all todo lists for the authenticated user, admins can see other users' todos too
func GetTodos(c *gin.Context) {
	username := getUsername(c)
	admin := getAdminState(c)
	var result []models.TodoList
	for _, todo := range models.MockTodoLists {
		if todo.DeletedAt != nil {
			continue
		}
		if todo.Username == username || admin {
			// Create a copy of the todo but only include non-deleted steps
			todoWithActiveSteps := todo
			var activeSteps []models.TodoStep
			for _, step := range todo.Steps {
				if step.DeletedAt == nil {
					activeSteps = append(activeSteps, step)
				}
			}
			todoWithActiveSteps.Steps = activeSteps
			todoWithActiveSteps.CalculateCompletion()
			result = append(result, todoWithActiveSteps)
		}
	}
	c.JSON(http.StatusOK, result)
}

// AddTodoList creates a new todo list for the authenticated user
func AddTodoList(c *gin.Context) {
	username := getUsername(c)
	var todo models.TodoList

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := todo.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	todo.ID = len(models.MockTodoLists) + 1
	todo.Username = username
	todo.CreatedAt = now
	todo.UpdatedAt = now
	todo.DeletedAt = nil

	for i := range todo.Steps {
		todo.Steps[i].ID = i + 1
		todo.Steps[i].TodoListID = todo.ID
		todo.Steps[i].CreatedAt = now
		todo.Steps[i].UpdatedAt = now
	}

	todo.CalculateCompletion()
	models.MockTodoLists = append(models.MockTodoLists, todo)
	c.JSON(http.StatusCreated, todo)
}

// RenameTodoList renames an existing todo list
func RenameTodoList(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	//admin := getAdminState(c)

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
				return
			}

			var update models.TodoList

			if err := c.ShouldBindJSON(&update); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := update.Validate(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			models.MockTodoLists[i].Name = update.Name
			models.MockTodoLists[i].UpdatedAt = time.Now()
			models.MockTodoLists[i].CalculateCompletion()

			c.JSON(http.StatusOK, models.MockTodoLists[i])
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found"})
}

// DeleteTodoList marks a todo list as deleted
func DeleteTodoList(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	//admin := getAdminState(c)
	now := time.Now()

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "It's already deleted"})
				return
			}
			models.MockTodoLists[i].DeletedAt = &now
			models.MockTodoLists[i].UpdatedAt = now
			c.JSON(http.StatusOK, gin.H{"message": "Todo list deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
}

// GetAllTodosAdmin returns all todo lists including deleted elements
// You can uncomment this function if you want the admin to be able to see all todos including deleted ones
/*
func GetAllTodosAdmin(c *gin.Context) {
	// Check if the user is admin
	admin := getAdminState(c)

	if !admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	// Return all todo lists including deleted elements
	c.JSON(http.StatusOK, models.MockTodoLists)
}
*/

// AddTodoStep adds a new step to an existing todo list
func AddTodoStep(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	//admin := getAdminState(c)

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
				return
			}

			var step models.TodoStep

			if err := c.ShouldBindJSON(&step); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			if err := step.Validate(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			now := time.Now()
			step.ID = len(todo.Steps) + 1
			step.TodoListID = todo.ID
			step.CreatedAt = now
			step.UpdatedAt = now
			step.DeletedAt = nil

			models.MockTodoLists[i].Steps = append(models.MockTodoLists[i].Steps, step)
			models.MockTodoLists[i].CalculateCompletion()

			c.JSON(http.StatusCreated, step)
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you don't have permission to that todo list"})
}

// DeleteTodoStep marks a step as deleted
func DeleteTodoStep(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	stepID, _ := strconv.Atoi(c.Param("step_id"))
	//admin := getAdminState(c)

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
				return
			}
			for j, step := range todo.Steps {
				if step.ID == stepID {
					if step.DeletedAt != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Step already deleted"})
						return
					}
					now := time.Now()
					models.MockTodoLists[i].Steps[j].DeletedAt = &now
					models.MockTodoLists[i].Steps[j].UpdatedAt = now
					models.MockTodoLists[i].CalculateCompletion()
					c.JSON(http.StatusOK, gin.H{"message": "Step deleted successfully"})
					return
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Step not found"})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you don't have permission to that todo list"})
}

// RenameTodoStep renames an existing step in a todo list
func RenameTodoStep(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	stepID, _ := strconv.Atoi(c.Param("step_id"))
	//admin := getAdminState(c)

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
				return
			}

			for j, step := range todo.Steps {
				if step.ID == stepID {
					if step.DeletedAt != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Step already deleted"})
						return
					}
					var update models.TodoStep

					if err := c.ShouldBindJSON(&update); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
						return
					}

					if err := update.Validate(); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}

					models.MockTodoLists[i].Steps[j].Content = update.Content
					models.MockTodoLists[i].Steps[j].UpdatedAt = time.Now()
					models.MockTodoLists[i].CalculateCompletion()

					c.JSON(http.StatusOK, models.MockTodoLists[i].Steps[j])
					return
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Step not found"})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you don't have permission to that todo list"})
}

// ToggleStepCompletion toggles the completion status of a step in a todo list
func ToggleStepCompletion(c *gin.Context) {
	username := getUsername(c)
	id, _ := strconv.Atoi(c.Param("id"))
	stepID, _ := strconv.Atoi(c.Param("step_id"))
	//admin := getAdminState(c)

	for i, todo := range models.MockTodoLists {
		if todo.ID == id && (todo.Username == username /*|| admin*/) {
			if todo.DeletedAt != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you have no permission to that todo list"})
				return
			}

			for j, step := range todo.Steps {
				if step.ID == stepID {
					if step.DeletedAt != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Step already deleted"})
						return
					}
					models.MockTodoLists[i].Steps[j].IsCompleted = !models.MockTodoLists[i].Steps[j].IsCompleted
					models.MockTodoLists[i].Steps[j].UpdatedAt = time.Now()
					models.MockTodoLists[i].CalculateCompletion()

					c.JSON(http.StatusOK, models.MockTodoLists[i].Steps[j])
					return
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Step not found"})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found or you don't have permission to that todo list"})
}
