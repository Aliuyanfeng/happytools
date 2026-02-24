// Package store
// @author: liuyanfeng
// @date: 2026/2/3
// @description: Category 分类增删改查操作
package store

import (
	"encoding/json"
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

// Category 分类
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllCategories 获取所有分类
func GetAllCategories() ([]Category, error) {
	categories := make([]Category, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(categoryBucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var c Category
			if err := json.Unmarshal(v, &c); err != nil {
				return err
			}
			categories = append(categories, c)
			return nil
		})
	})
	return categories, err
}

// GetCategory 根据 ID 获取分类
func GetCategory(id int) (*Category, error) {
	var category Category
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(categoryBucket)
		if b == nil {
			return errors.New("category not found")
		}
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("category not found")
		}
		return json.Unmarshal(v, &category)
	})
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateCategory 创建分类
func CreateCategory(name, color string) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	var category Category
	err := DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(categoryBucket)

		// 获取自增 ID
		id, _ := b.NextSequence()
		now := time.Now()

		// 计算下一个排序顺序
		sortOrder := 0
		b.ForEach(func(k, v []byte) error {
			var c Category
			if json.Unmarshal(v, &c) == nil && c.SortOrder >= sortOrder {
				sortOrder = c.SortOrder + 1
			}
			return nil
		})

		category = Category{
			ID:        int(id),
			Name:      name,
			Color:     color,
			SortOrder: sortOrder,
			CreatedAt: now,
			UpdatedAt: now,
		}

		data, err := json.Marshal(category)
		if err != nil {
			return err
		}
		return b.Put(itob(category.ID), data)
	})
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory 更新分类
func UpdateCategory(id int, name, color string) error {
	if name == "" {
		return errors.New("category name cannot be empty")
	}

	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(categoryBucket)
		v := b.Get(itob(id))
		if v == nil {
			return errors.New("category not found")
		}

		var category Category
		if err := json.Unmarshal(v, &category); err != nil {
			return err
		}

		category.Name = name
		category.Color = color
		category.UpdatedAt = time.Now()

		data, err := json.Marshal(category)
		if err != nil {
			return err
		}
		return b.Put(itob(id), data)
	})
}

// DeleteCategory 删除分类
func DeleteCategory(id int) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		// 检查是否有待办关联此分类
		todoBkt := tx.Bucket(todoBucket)
		if todoBkt != nil {
			err := todoBkt.ForEach(func(k, v []byte) error {
				var todo Todo
				if err := json.Unmarshal(v, &todo); err == nil {
					if todo.CategoryID != nil && *todo.CategoryID == id {
						return errors.New("cannot delete category: todos are associated with it")
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		b := tx.Bucket(categoryBucket)
		return b.Delete(itob(id))
	})
}
