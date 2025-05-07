package models

import (
	"errors"
	"strings"
	"time"
)

// TO-DO List entity
type TodoList struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Completion float64    `json:"completion"` // % olarak
	Username   string     `json:"username"`
	Steps      []TodoStep `json:"steps"`
}

// TO-DO Step entity - steps of a todo list
type TodoStep struct {
	ID          int        `json:"id"`
	TodoListID  int        `json:"todo_list_id"` // Hangi listeye ait
	Content     string     `json:"content"`
	IsCompleted bool       `json:"is_completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (t *TodoList) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("list name cannot be empty")
	}
	if t.Completion < 0 || t.Completion > 100 {
		return errors.New("completion must be between 0 and 100")
	}
	return nil
}

func (s *TodoStep) Validate() error {
	if strings.TrimSpace(s.Content) == "" {
		return errors.New("step content cannot be empty")
	}
	return nil
}

// CalculateCompletion updates the completion percentage based on completed steps
func (t *TodoList) CalculateCompletion() {
	if len(t.Steps) == 0 {
		t.Completion = 0
		return
	}

	completedCount := 0
	for _, step := range t.Steps {
		if step.IsCompleted {
			completedCount++
		}
	}

	t.Completion = float64(completedCount) / float64(len(t.Steps)) * 100
}

// Mock todo lists for testing purposes
var MockTodoLists = []TodoList{
	{
		ID:         1,
		Name:       "Günlük Rutin",
		CreatedAt:  time.Now().AddDate(0, 0, -2),
		UpdatedAt:  time.Now(),
		Completion: 66.7,
		Username:   "admin",
		Steps: []TodoStep{
			{
				ID:          1,
				TodoListID:  1,
				Content:     "Kahvaltı yap",
				IsCompleted: true,
				CreatedAt:   time.Now().AddDate(0, 0, -2),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				TodoListID:  1,
				Content:     "Spor yap",
				IsCompleted: true,
				CreatedAt:   time.Now().AddDate(0, 0, -2),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          3,
				TodoListID:  1,
				Content:     "Kitap oku",
				IsCompleted: false,
				CreatedAt:   time.Now().AddDate(0, 0, -2),
				UpdatedAt:   time.Now(),
			},
		},
	},
	{
		ID:         2,
		Name:       "Projeye Hazırlık",
		CreatedAt:  time.Now().AddDate(0, 0, -5),
		UpdatedAt:  time.Now(),
		Completion: 33.3,
		Username:   "user",
		Steps: []TodoStep{
			{
				ID:          1,
				TodoListID:  2,
				Content:     "Projeyi araştır",
				IsCompleted: true,
				CreatedAt:   time.Now().AddDate(0, 0, -5),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				TodoListID:  2,
				Content:     "Gereksinimleri yaz",
				IsCompleted: false,
				CreatedAt:   time.Now().AddDate(0, 0, -4),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          3,
				TodoListID:  2,
				Content:     "Tasarım yap",
				IsCompleted: false,
				CreatedAt:   time.Now().AddDate(0, 0, -3),
				UpdatedAt:   time.Now(),
			},
		},
	},
}
