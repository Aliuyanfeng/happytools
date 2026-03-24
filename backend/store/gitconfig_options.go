// Package store
// @description: Git Config Manager 存储层 — 仓库记录 & QuickPanel 配置
package store

import (
	"encoding/json"
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

// Repository 受管理的 Git 仓库记录
type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Platform  string    `json:"platform"` // "github" | "gitlab" | "gitee" | "custom"
	CreatedAt time.Time `json:"createdAt"`
}

// QuickPanelItem QuickPanel 中的一个配置项
type QuickPanelItem struct {
	Section string `json:"section"`
	SubKey  string `json:"subKey"`
	Key     string `json:"key"`
	Order   int    `json:"order"`
}

// SaveRepository 新增或更新仓库记录
func SaveRepository(repo *Repository) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitRepoBucket)
		data, err := json.Marshal(repo)
		if err != nil {
			return err
		}
		return b.Put([]byte(repo.ID), data)
	})
}

// GetAllRepositories 获取所有仓库记录
func GetAllRepositories() ([]*Repository, error) {
	repos := make([]*Repository, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitRepoBucket)
		return b.ForEach(func(k, v []byte) error {
			var r Repository
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			repos = append(repos, &r)
			return nil
		})
	})
	return repos, err
}

// GetRepository 根据 ID 获取仓库记录
func GetRepository(id string) (*Repository, error) {
	var repo Repository
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitRepoBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("repository not found")
		}
		return json.Unmarshal(v, &repo)
	})
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// DeleteRepository 删除仓库记录
func DeleteRepository(id string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitRepoBucket)
		return b.Delete([]byte(id))
	})
}

// GetQuickPanel 获取指定仓库的 QuickPanel 配置
func GetQuickPanel(repoID string) ([]QuickPanelItem, error) {
	var items []QuickPanelItem
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitQuickPanelBucket)
		v := b.Get([]byte(repoID))
		if v == nil {
			// 返回默认三项
			items = defaultQuickPanel()
			return nil
		}
		return json.Unmarshal(v, &items)
	})
	return items, err
}

// SaveQuickPanel 保存指定仓库的 QuickPanel 配置
func SaveQuickPanel(repoID string, items []QuickPanelItem) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(gitQuickPanelBucket)
		data, err := json.Marshal(items)
		if err != nil {
			return err
		}
		return b.Put([]byte(repoID), data)
	})
}

// defaultQuickPanel 返回默认 QuickPanel 配置
func defaultQuickPanel() []QuickPanelItem {
	return []QuickPanelItem{
		{Section: "user", SubKey: "", Key: "name", Order: 0},
		{Section: "user", SubKey: "", Key: "email", Order: 1},
		{Section: "remote", SubKey: "origin", Key: "url", Order: 2},
	}
}
