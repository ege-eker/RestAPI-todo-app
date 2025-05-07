# RESTful Todo API

A simple RESTful API for managing todo lists using Go and the Gin framework.

## Table of Contents

- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Models](#models)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Usage Examples](#usage-examples)

## Features

- User authentication with JWT (JSON Web Tokens)
- CRUD operations for todo lists and their steps
- Role-based access control (regular users and admins)
- Soft deletion for todo lists and steps
- Automatic completion percentage calculation

## API Endpoints

### Authentication
- `POST /api/login` - Authenticate user and get JWT token

### Todo Lists (Requires Authentication)
- `GET /api/todo-lists` - Get all todo lists for the authenticated user
- `POST /api/todo-lists` - Create a new todo list
- `PUT /api/todo-lists/:id` - Rename a todo list
- `DELETE /api/todo-lists/:id` - Delete a todo list (soft delete)

### Todo Steps (Requires Authentication)
- `POST /api/todo-lists/:id/step` - Add a step to a todo list
- `PUT /api/todo-lists/:id/step/:step_id` - Rename a step
- `DELETE /api/todo-lists/:id/step/:step_id` - Delete a step (soft delete)
- `POST /api/todo-lists/:id/step/:step_id/toggle` - Toggle completion status of a step

### Admin Endpoints (Optional)
- `GET /api/admin/todolists` - Get all todo lists including deleted elements (admin only, uncomment in main.go to enable)

## Authentication

The API uses JWT for authentication. Upon successful login, a token is returned which must be included in the `Authorization` header of subsequent requests as a Bearer token:

```
Authorization: Bearer <token>
```

### Demo Users
- Admin: username: `admin`, password: `admin123`
- Regular User: username: `user`, password: `password123`

## Models

### TodoList
```go
type TodoList struct {
    ID         int        `json:"id"`
    Name       string     `json:"name"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at,omitempty"`
    Completion float64    `json:"completion"` // Percentage
    Username   string     `json:"username"`
    Steps      []TodoStep `json:"steps"`
}
```

### TodoStep
```go
type TodoStep struct {
    ID          int        `json:"id"`
    TodoListID  int        `json:"todo_list_id"`
    Content     string     `json:"content"`
    IsCompleted bool       `json:"is_completed"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
```

## Getting Started

### Prerequisites

- Go (1.16 or higher)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/RestAPI-todo-app.git
   cd RestAPI-todo-app
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Application

Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`.

## Usage Examples

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "password123"}'
```

### Get Todo Lists
```bash
curl -X GET http://localhost:8080/api/todo-lists \
  -H "Authorization: Bearer <your_token>"
```

### Create a Todo List
```bash
curl -X POST http://localhost:8080/api/todo-lists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{"name": "Shopping List", "steps": [{"content": "Buy milk"}, {"content": "Buy eggs"}]}'
```

### Add a Step to a Todo List
```bash
curl -X POST http://localhost:8080/api/todo-lists/1/step \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{"content": "New task"}'
```

### Toggle Step Completion
```bash
curl -X POST http://localhost:8080/api/todo-lists/1/step/1/toggle \
  -H "Authorization: Bearer <your_token>"
```

## Security Notes

This project uses a hardcoded JWT secret for demonstration purposes only. In a production environment, you should:

1. Use environment variables or a secure configuration system for the JWT secret
2. Implement HTTPS to secure the communication
3. Add proper input validation and sanitization
4. Consider implementing rate limiting
