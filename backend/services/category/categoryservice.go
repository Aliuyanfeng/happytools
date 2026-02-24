// Package category
// @author: liuyanfeng
// @date: 2026/2/3
// @description: 分类服务
package category

import (
	"errors"

	"github.com/Aliuyanfeng/happytools/backend/store"
)

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
}

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// GetAll 获取所有分类
func (c *CategoryService) GetAll() ([]Category, error) {
	categories, err := store.GetAllCategories()
	if err != nil {
		return nil, err
	}

	result := make([]Category, len(categories))
	for i, cat := range categories {
		result[i] = Category{
			ID:        cat.ID,
			Name:      cat.Name,
			Color:     cat.Color,
			SortOrder: cat.SortOrder,
		}
	}
	return result, nil
}

// Add 添加分类
func (c *CategoryService) Add(name, color string) (*Category, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	cat, err := store.CreateCategory(name, color)
	if err != nil {
		return nil, err
	}

	return &Category{
		ID:        cat.ID,
		Name:      cat.Name,
		Color:     cat.Color,
		SortOrder: cat.SortOrder,
	}, nil
}

// Update 更新分类
func (c *CategoryService) Update(id int, name, color string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return store.UpdateCategory(id, name, color)
}

// Delete 删除分类
func (c *CategoryService) Delete(id int) error {
	return store.DeleteCategory(id)
}
