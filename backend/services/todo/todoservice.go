// Package todo
// @author: liuyanfeng
// @date: 2026/1/9 14:37
// @description:
package todo

import (
	"errors"
	"time"

	"github.com/Aliuyanfeng/happytools/backend/store"
)

type Todo struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Completed  bool    `json:"completed"`
	CategoryID *int    `json:"category_id"`
	DueDate    *string `json:"due_date"`
	Priority   int     `json:"priority"`
	Status     int     `json:"status"`
}

type TodoService struct{}

func NewTodoService() *TodoService {
	return &TodoService{}
}

func (t *TodoService) GetAll() ([]Todo, error) {
	todos, err := store.GetAllTodos()
	if err != nil {
		return nil, err
	}

	result := make([]Todo, len(todos))
	for i, item := range todos {
		var dueDateStr *string
		if item.DueDate != nil {
			s := item.DueDate.Format("2006-01-02")
			dueDateStr = &s
		}
		result[i] = Todo{
			ID:         item.ID,
			Title:      item.Title,
			Completed:  item.Completed,
			CategoryID: item.CategoryID,
			DueDate:    dueDateStr,
			Priority:   item.Priority,
			Status:     item.GetStatus(),
		}
	}
	return result, nil
}

func (t *TodoService) Add(title string, categoryID *int, dueDate *string, priority int) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	var parsedDue *time.Time
	if dueDate != nil && *dueDate != "" {
		d, err := time.Parse("2006-01-02", *dueDate)
		if err != nil {
			return nil, err
		}
		parsedDue = &d
	}

	todo, err := store.CreateTodoEnhanced(title, categoryID, parsedDue, priority)
	if err != nil {
		return nil, err
	}

	var dueDateStr *string
	if todo.DueDate != nil {
		s := todo.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	return &Todo{
		ID:         todo.ID,
		Title:      todo.Title,
		Completed:  todo.Completed,
		CategoryID: todo.CategoryID,
		DueDate:    dueDateStr,
		Priority:   todo.Priority,
		Status:     todo.GetStatus(),
	}, nil
}

func (t *TodoService) Update(id int, title string, completed bool, categoryID *int, dueDate *string, priority int) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	var parsedDue *time.Time
	if dueDate != nil && *dueDate != "" {
		d, err := time.Parse("2006-01-02", *dueDate)
		if err != nil {
			return err
		}
		parsedDue = &d
	}

	return store.UpdateTodoEnhanced(id, title, completed, categoryID, parsedDue, priority)
}

func (t *TodoService) Toggle(id int) error {
	return store.ToggleTodo(id)
}

func (t *TodoService) Delete(id int) error {
	return store.DeleteTodo(id)
}
