// Package virusTotal
// @author: liuyanfeng
// @date: 2026/3/2
// @description: VirusTotal 任务管理
package virusTotal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Aliuyanfeng/happytools/backend/store"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks map[string]*VTTask
	mu    sync.RWMutex
}

var taskManager *TaskManager
var taskManagerOnce sync.Once

// GetTaskManager 获取任务管理器单例
func GetTaskManager() *TaskManager {
	taskManagerOnce.Do(func() {
		taskManager = &TaskManager{
			tasks: make(map[string]*VTTask),
		}
		taskManager.loadTasksFromFile()
	})
	return taskManager
}

// saveTasksToFile 保存任务到文件（备份）
func (tm *TaskManager) saveTasksToFile() {
	dataDir := filepath.Join(os.Getenv("APPDATA"), "HappyTools", "vt")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Failed to create data directory: %v\n", err)
		return
	}

	dataFile := filepath.Join(dataDir, "tasks.json")
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal tasks: %v\n", err)
		return
	}

	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		fmt.Printf("Failed to save tasks: %v\n", err)
	}
}

// loadTasksFromFile 从文件加载任务（备份）
func (tm *TaskManager) loadTasksFromFile() {
	dataDir := filepath.Join(os.Getenv("APPDATA"), "HappyTools", "vt")
	dataFile := filepath.Join(dataDir, "tasks.json")

	data, err := os.ReadFile(dataFile)
	if err != nil {
		// 文件不存在，创建空任务列表
		return
	}

	if err := json.Unmarshal(data, &tm.tasks); err != nil {
		fmt.Printf("Failed to unmarshal tasks: %v\n", err)
	}
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("VT-%d", time.Now().UnixNano())
}

// generateFileID 生成文件ID
func generateFileID() string {
	return fmt.Sprintf("FILE-%d", time.Now().UnixNano())
}

// SaveTask 保存任务到数据库
func (tm *TaskManager) SaveTask(task *VTTask) error {
	tm.mu.Lock()
	tm.tasks[task.ID] = task
	tm.mu.Unlock()

	// 转换为 store 的任务格式
	storeTask := &store.VTTask{
		ID:              task.ID,
		TaskType:        store.TaskType(task.TaskType),
		FileName:        task.FileName,
		FilePath:        task.FilePath,
		FileSize:        task.FileSize,
		Status:          task.Status,
		AnalysisID:      task.AnalysisID,
		CreatedAt:       task.CreatedAt,
		CompletedAt:     task.CompletedAt,
		MD5:             task.MD5,
		SHA256:          task.SHA256,
		SHA1:            task.SHA1,
		DetectionRate:   task.DetectionRate,
		TotalEngines:    task.TotalEngines,
		ScanTime:        task.ScanTime,
		Stats:           store.StatsInfo(task.Stats),
		TotalFiles:      task.TotalFiles,
		MaliciousFiles:  task.MaliciousFiles,
		SuspiciousFiles: task.SuspiciousFiles,
		CompletedFiles:  task.CompletedFiles,
	}

	// 同时保存到文件作为备份
	tm.saveTasksToFile()

	return store.SaveVTTask(storeTask)
}

// LoadTask 从数据库加载任务
func (tm *TaskManager) LoadTask(id string) (*VTTask, error) {
	// 先从内存缓存获取
	tm.mu.RLock()
	if task, exists := tm.tasks[id]; exists {
		tm.mu.RUnlock()
		return task, nil
	}
	tm.mu.RUnlock()

	storeTask, err := store.GetVTTask(id)
	if err != nil {
		return nil, err
	}

	task := &VTTask{
		ID:              storeTask.ID,
		TaskType:        string(storeTask.TaskType),
		FileName:        storeTask.FileName,
		FilePath:        storeTask.FilePath,
		FileSize:        storeTask.FileSize,
		Status:          storeTask.Status,
		AnalysisID:      storeTask.AnalysisID,
		CreatedAt:       storeTask.CreatedAt,
		CompletedAt:     storeTask.CompletedAt,
		MD5:             storeTask.MD5,
		SHA256:          storeTask.SHA256,
		SHA1:            storeTask.SHA1,
		DetectionRate:   storeTask.DetectionRate,
		TotalEngines:    storeTask.TotalEngines,
		ScanTime:        storeTask.ScanTime,
		Stats:           StatsInfo(storeTask.Stats),
		TotalFiles:      storeTask.TotalFiles,
		MaliciousFiles:  storeTask.MaliciousFiles,
		SuspiciousFiles: storeTask.SuspiciousFiles,
		CompletedFiles:  storeTask.CompletedFiles,
	}

	// 缓存到内存
	tm.mu.Lock()
	tm.tasks[id] = task
	tm.mu.Unlock()

	return task, nil
}

// LoadAllTasks 从数据库加载所有任务
func (tm *TaskManager) LoadAllTasks() ([]*VTTask, error) {
	storeTasks, err := store.GetAllVTTasks()
	if err != nil {
		return nil, err
	}

	tasks := make([]*VTTask, len(storeTasks))
	for i, storeTask := range storeTasks {
		tasks[i] = &VTTask{
			ID:              storeTask.ID,
			TaskType:        string(storeTask.TaskType),
			FileName:        storeTask.FileName,
			FilePath:        storeTask.FilePath,
			FileSize:        storeTask.FileSize,
			Status:          storeTask.Status,
			AnalysisID:      storeTask.AnalysisID,
			CreatedAt:       storeTask.CreatedAt,
			CompletedAt:     storeTask.CompletedAt,
			MD5:             storeTask.MD5,
			SHA256:          storeTask.SHA256,
			SHA1:            storeTask.SHA1,
			DetectionRate:   storeTask.DetectionRate,
			TotalEngines:    storeTask.TotalEngines,
			ScanTime:        storeTask.ScanTime,
			Stats:           StatsInfo(storeTask.Stats),
			TotalFiles:      storeTask.TotalFiles,
			MaliciousFiles:  storeTask.MaliciousFiles,
			SuspiciousFiles: storeTask.SuspiciousFiles,
			CompletedFiles:  storeTask.CompletedFiles,
		}
	}

	return tasks, nil
}

// DeleteTask 删除任务
func (tm *TaskManager) DeleteTask(id string) error {
	tm.mu.Lock()
	delete(tm.tasks, id)
	tm.mu.Unlock()

	tm.saveTasksToFile()

	// 同时删除关联的子文件
	store.DeleteVTFilesByTaskID(id)

	return store.DeleteVTTask(id)
}
