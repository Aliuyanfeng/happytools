// Package store
// @author: liuyanfeng
// @date: 2026/3/2
// @description: VirusTotal 任务存储
package store

import (
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeSingle    TaskType = "single"    // 单文件扫描
	TaskTypeDirectory TaskType = "directory" // 目录批量扫描
)

// VTTask VirusTotal 任务
type VTTask struct {
	ID            string    `json:"id"`
	TaskType      TaskType  `json:"taskType"`      // 任务类型: single/directory
	FileName      string    `json:"fileName"`      // 单文件时为文件名，目录时为目录名
	FilePath      string    `json:"filePath"`      // 单文件时为文件路径，目录时为目录路径
	FileSize      int64     `json:"fileSize"`      // 单文件时为文件大小，目录时为总大小
	Status        string    `json:"status"`        // 任务状态
	AnalysisID    string    `json:"analysisId"`    // 单文件时的分析ID
	CreatedAt     time.Time `json:"createdAt"`
	CompletedAt   time.Time `json:"completedAt"`
	MD5           string    `json:"md5"`
	SHA256        string    `json:"sha256"`
	SHA1          string    `json:"sha1"`
	DetectionRate int       `json:"detectionRate"`
	TotalEngines  int       `json:"totalEngines"`
	ScanTime      string    `json:"scanTime"`
	Stats         StatsInfo `json:"stats"` // 单文件任务的检测统计

	// 目录扫描专用字段
	TotalFiles      int `json:"totalFiles"`      // 目录下文件总数
	MaliciousFiles  int `json:"maliciousFiles"`  // 带毒文件数
	SuspiciousFiles int `json:"suspiciousFiles"` // 可疑文件数
	CompletedFiles  int `json:"completedFiles"`  // 已完成扫描文件数
}

// VTFile 批量任务中的子文件
type VTFile struct {
	ID            string       `json:"id"`
	TaskID        string       `json:"taskId"`        // 所属任务ID
	FileName      string       `json:"fileName"`
	FilePath      string       `json:"filePath"`
	FileSize      int64        `json:"fileSize"`
	Status        string       `json:"status"`        // pending/queued/completed/failed
	AnalysisID    string       `json:"analysisId"`
	CreatedAt     time.Time    `json:"createdAt"`
	CompletedAt   time.Time    `json:"completedAt"`
	MD5           string       `json:"md5"`
	SHA256        string       `json:"sha256"`
	SHA1          string       `json:"sha1"`
	DetectionRate int          `json:"detectionRate"`
	TotalEngines  int          `json:"totalEngines"`
	Stats         StatsInfo    `json:"stats"`
	Results       []ScanResult `json:"results"` // 检测结果
}

// ScanResult 检测结果
type ScanResult struct {
	Engine   string `json:"engine"`
	Category string `json:"category"`
	Result   string `json:"result"`
	Method   string `json:"method"`
}

// StatsInfo 统计信息
type StatsInfo struct {
	Malicious        int `json:"malicious"`
	Suspicious       int `json:"suspicious"`
	Harmless         int `json:"harmless"`
	Undetected       int `json:"undetected"`
	TypeUnsupported  int `json:"typeUnsupported"`
	ConfirmedTimeout int `json:"confirmedTimeout"`
	Timeout          int `json:"timeout"`
	Failure          int `json:"failure"`
}

// SaveVTTask 保存任务
func SaveVTTask(task *VTTask) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(vtTaskBucket)
		if err != nil {
			return err
		}

		data, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(task.ID), data)
	})
}

// GetVTTask 获取任务
func GetVTTask(id string) (*VTTask, error) {
	var task VTTask

	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtTaskBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return bbolt.ErrBucketNotFound
		}

		return json.Unmarshal(data, &task)
	})

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// GetAllVTTasks 获取所有任务
func GetAllVTTasks() ([]*VTTask, error) {
	var tasks []*VTTask

	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtTaskBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var task VTTask
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, &task)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// DeleteVTTask 删除任务
func DeleteVTTask(id string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtTaskBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		return bucket.Delete([]byte(id))
	})
}

// ========== 子文件存储操作 ==========

// SaveVTFile 保存子文件
func SaveVTFile(file *VTFile) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(vtFileBucket)
		if err != nil {
			return err
		}

		data, err := json.Marshal(file)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(file.ID), data)
	})
}

// GetVTFile 获取子文件
func GetVTFile(id string) (*VTFile, error) {
	var file VTFile

	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtFileBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return bbolt.ErrBucketNotFound
		}

		return json.Unmarshal(data, &file)
	})

	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetVTFilesByTaskID 获取任务下所有子文件
func GetVTFilesByTaskID(taskID string) ([]*VTFile, error) {
	var files []*VTFile

	err := DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtFileBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var file VTFile
			if err := json.Unmarshal(v, &file); err != nil {
				return err
			}
			if file.TaskID == taskID {
				files = append(files, &file)
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// DeleteVTFile 删除子文件
func DeleteVTFile(id string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtFileBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		return bucket.Delete([]byte(id))
	})
}

// DeleteVTFilesByTaskID 删除任务下所有子文件
func DeleteVTFilesByTaskID(taskID string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(vtFileBucket)
		if bucket == nil {
			return nil
		}

		var keysToDelete [][]byte
		bucket.ForEach(func(k, v []byte) error {
			var file VTFile
			if err := json.Unmarshal(v, &file); err != nil {
				return err
			}
			if file.TaskID == taskID {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		for _, key := range keysToDelete {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}
