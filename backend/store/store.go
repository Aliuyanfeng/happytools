// Package store
// @author: liuyanfeng
// @date: 2026/2/3
// @description: bbolt 存储层基础初始化
package store

import (
	"encoding/binary"
	"time"

	"go.etcd.io/bbolt"
)

var (
	DB                      *bbolt.DB
	todoBucket              = []byte("todos")
	categoryBucket          = []byte("categories")
	dailyReportBucket       = []byte("daily_reports")
	appSettingsBucket       = []byte("app_settings")
	vtTaskBucket            = []byte("vt_tasks")
	vtFileBucket            = []byte("vt_files") // 批量任务中的子文件
	gitRepoBucket           = []byte("git_repos")
	gitQuickPanelBucket     = []byte("git_quick_panels")
	makefileRecentBucket    = []byte("makefile_recent")
	makefileTemplateBucket  = []byte("makefile_templates")
)

// Init 初始化 bbolt 数据库
func Init(path string) error {
	var err error
	DB, err = bbolt.Open(path, 0600, &bbolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}

	// 创建所有 bucket
	return DB.Update(func(tx *bbolt.Tx) error {
		buckets := [][]byte{todoBucket, categoryBucket, dailyReportBucket, appSettingsBucket, vtTaskBucket, vtFileBucket, gitRepoBucket, gitQuickPanelBucket, makefileRecentBucket, makefileTemplateBucket}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

// Close 关闭数据库
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// itob 将 int 转换为 []byte (用于 key)
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
