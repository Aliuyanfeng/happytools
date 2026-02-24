// Package store
// @author: liuyanfeng
// @date: 2026/2/6
// @description: 应用设置存储
package store

import (
	"time"

	"go.etcd.io/bbolt"
)

const (
	LastUsedTimeKey = "last_used_time"
)

// GetLastUsedTime 获取上次使用时间
func GetLastUsedTime() (*time.Time, error) {
	var lastUsed *time.Time
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(appSettingsBucket)
		if b == nil {
			return nil
		}

		data := b.Get([]byte(LastUsedTimeKey))
		if data == nil {
			return nil
		}

		t, err := time.Parse(time.RFC3339, string(data))
		if err != nil {
			return err
		}

		lastUsed = &t
		return nil
	})

	return lastUsed, err
}

// UpdateLastUsedTime 更新上次使用时间为当前时间
func UpdateLastUsedTime() error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(appSettingsBucket)
		if b == nil {
			return nil
		}

		now := time.Now().Format(time.RFC3339)
		return b.Put([]byte(LastUsedTimeKey), []byte(now))
	})
}
