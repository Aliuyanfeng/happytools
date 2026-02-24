// Package store
// @author: liuyanfeng
// @date: 2026/2/3
// @description: Todo 增删改查操作
package store

import (
	"encoding/json"
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

// Todo 待办事项
type Todo struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Completed  bool       `json:"completed"`
	CategoryID *int       `json:"category_id"`
	DueDate    *time.Time `json:"due_date"`
	Priority   int        `json:"priority"` // 0-低 1-中 2-高
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TodoStatus 任务状态常量
const (
	TodoStatusNormal  = 0 // 正常
	TodoStatusWarning = 1 // 即将到期 (24小时内)
	TodoStatusOverdue = 2 // 已逾期
)

// GetStatus 获取任务状态
func (t *Todo) GetStatus() int {
	if t.DueDate == nil || t.Completed {
		return TodoStatusNormal
	}

	now := time.Now()
	if t.DueDate.Before(now) {
		return TodoStatusOverdue
	}

	// 24小时内到期视为警告
	if t.DueDate.Sub(now) <= 24*time.Hour {
		return TodoStatusWarning
	}

	return TodoStatusNormal
}

// GetAllTodos 获取所有待办
func GetAllTodos() ([]Todo, error) {
	todos := make([]Todo, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		return b.ForEach(func(k, v []byte) error {
			var t Todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			todos = append(todos, t)
			return nil
		})
	})
	return todos, err
}

// GetTodo 根据 ID 获取待办
func GetTodo(id int) (*Todo, error) {
	var todo Todo
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("todo not found")
		}
		return json.Unmarshal(v, &todo)
	})
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// CreateTodo 创建待办
func CreateTodo(title string) (*Todo, error) {
	var todo Todo
	err := DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)

		// 获取自增 ID
		id, _ := b.NextSequence()
		now := time.Now()

		todo = Todo{
			ID:        int(id),
			Title:     title,
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(itob(todo.ID), data)
	})
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo 更新待办
func UpdateTodo(id int, title string, completed bool) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("todo not found")
		}

		var todo Todo
		if err := json.Unmarshal(v, &todo); err != nil {
			return err
		}

		todo.Title = title
		todo.Completed = completed
		todo.UpdatedAt = time.Now()

		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(itob(id), data)
	})
}

// ToggleTodo 切换待办完成状态
func ToggleTodo(id int) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("todo not found")
		}

		var todo Todo
		if err := json.Unmarshal(v, &todo); err != nil {
			return err
		}

		todo.Completed = !todo.Completed
		todo.UpdatedAt = time.Now()

		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(itob(id), data)
	})
}

// DeleteTodo 删除待办
func DeleteTodo(id int) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		return b.Delete(itob(id))
	})
}

// CreateTodoEnhanced 创建待办 (增强版)
func CreateTodoEnhanced(title string, categoryID *int, dueDate *time.Time, priority int) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	var todo Todo
	err := DB.Update(func(tx *bbolt.Tx) error {
		// 验证分类是否存在
		if categoryID != nil {
			catBkt := tx.Bucket(categoryBucket)
			if catBkt != nil && catBkt.Get(itob(*categoryID)) == nil {
				return errors.New("category not found")
			}
		}

		b := tx.Bucket(todoBucket)
		id, _ := b.NextSequence()
		now := time.Now()

		todo = Todo{
			ID:         int(id),
			Title:      title,
			Completed:  false,
			CategoryID: categoryID,
			DueDate:    dueDate,
			Priority:   priority,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(itob(todo.ID), data)
	})
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodoEnhanced 更新待办 (增强版)
func UpdateTodoEnhanced(id int, title string, completed bool, categoryID *int, dueDate *time.Time, priority int) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	return DB.Update(func(tx *bbolt.Tx) error {
		// 验证分类是否存在
		if categoryID != nil {
			catBkt := tx.Bucket(categoryBucket)
			if catBkt != nil && catBkt.Get(itob(*categoryID)) == nil {
				return errors.New("category not found")
			}
		}

		b := tx.Bucket(todoBucket)
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("todo not found")
		}

		var todo Todo
		if err := json.Unmarshal(v, &todo); err != nil {
			return err
		}

		todo.Title = title
		todo.Completed = completed
		todo.CategoryID = categoryID
		todo.DueDate = dueDate
		todo.Priority = priority
		todo.UpdatedAt = time.Now()

		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(itob(id), data)
	})
}

// GetTodosByCategory 按分类获取待办
func GetTodosByCategory(categoryID int) ([]Todo, error) {
	todos := make([]Todo, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		return b.ForEach(func(k, v []byte) error {
			var t Todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			if t.CategoryID != nil && *t.CategoryID == categoryID {
				todos = append(todos, t)
			}
			return nil
		})
	})
	return todos, err
}

// GetOverdueTodos 获取已逾期的待办
func GetOverdueTodos() ([]Todo, error) {
	todos := make([]Todo, 0)
	now := time.Now()

	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(todoBucket)
		return b.ForEach(func(k, v []byte) error {
			var t Todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			if !t.Completed && t.DueDate != nil && t.DueDate.Before(now) {
				todos = append(todos, t)
			}
			return nil
		})
	})
	return todos, err
}
