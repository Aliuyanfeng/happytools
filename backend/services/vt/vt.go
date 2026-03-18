// Package virusTotal
// @author: liuyanfeng
// @date: 2025/9/5 17:48
// @description: VT service for file operations
package virusTotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"log"
	"github.com/Aliuyanfeng/happytools/backend/store"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// 全局 VTService 实例
var (
	vtServiceInstance *VTService
	vtServiceOnce     sync.Once
)

// VTService provides file-related operations.
type VTService struct {
	apiKey       string
	app          *application.App
	concurrency  int // 并发扫描数
	scanQueue    chan *scanJob
	scanWg       sync.WaitGroup
	scanStopChan chan struct{}
}

// scanJob 扫描任务
type scanJob struct {
	file   *store.VTFile
	taskID string
}

// NewVTService creates a new instance of VTService.
func NewVTService(app *application.App) *VTService {
	vtServiceOnce.Do(func() {
		vtServiceInstance = &VTService{
			app:          app,
			concurrency:  5, // 默认并发数
			scanQueue:    make(chan *scanJob, 100),
			scanStopChan: make(chan struct{}),
		}
		// 启动扫描工作池
		vtServiceInstance.startScanWorkers()
	})
	return vtServiceInstance
}

// GetVTService 获取全局 VTService 实例
func GetVTService() *VTService {
	return vtServiceInstance
}

// UserInfo VirusTotal 用户基础信息
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Type      string `json:"type"`
	LastLogin int64  `json:"lastLogin"`
}

// GetUserInfo 获取当前 API Key 对应的用户信息
func (v *VTService) GetUserInfo(apiKey string) (*UserInfo, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.virustotal.com/api/v3/users/%s", apiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("x-apikey", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Email     string `json:"email"`
				LastLogin int64  `json:"last_login"`
			} `json:"attributes"`
		} `json:"data"`
	}
	log.Printf(fmt.Sprintf("%+v", json.NewDecoder(resp.Body)))
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &UserInfo{
		ID:        result.Data.ID,
		Type:      result.Data.Type,
		Email:     result.Data.Attributes.Email,
		LastLogin: result.Data.Attributes.LastLogin,
	}, nil
}

// SetAPIKey sets the API key for VirusTotal
func (v *VTService) SetAPIKey(apiKey string) error {
	fmt.Printf("当前API-KEY=%s", apiKey)
	v.apiKey = apiKey
	return nil
}

// GetAPIKey 获取当前 API Key
func (v *VTService) GetAPIKey() string {
	return v.apiKey
}

// SetConcurrency 设置并发扫描数
func (v *VTService) SetConcurrency(concurrency int) error {
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > 10 {
		concurrency = 10
	}
	v.concurrency = concurrency
	return nil
}

// GetConcurrency 获取并发扫描数
func (v *VTService) GetConcurrency() int {
	return v.concurrency
}

// ScanResponse represents the response from VirusTotal scan API
type ScanResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
}

// AnalysisResponse represents the response from VirusTotal analysis API
type AnalysisResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status  string `json:"status"`
			Results map[string]struct {
				Category string `json:"category"`
				Result   string `json:"result"`
				Method   string `json:"method"`
				EngineID string `json:"engine_id"`
			} `json:"results"`
			Stats struct {
				Harmless         int `json:"harmless"`
				TypeUnsupported  int `json:"type-unsupported"`
				Suspicious       int `json:"suspicious"`
				ConfirmedTimeout int `json:"confirmed-timeout"`
				Timeout          int `json:"timeout"`
				Failure          int `json:"failure"`
				Malicious        int `json:"malicious"`
				Undetected       int `json:"undetected"`
			} `json:"stats"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		FileInfo struct {
			Md5    string `json:"md5"`
			Sha1   string `json:"sha1"`
			Sha256 string `json:"sha256"`
			Size   int    `json:"size"`
		} `json:"file_info"`
	} `json:"meta"`
}

// ScanResult represents a single scan result
type ScanResult struct {
	Engine   string `json:"engine"`
	Category string `json:"category"`
	Result   string `json:"result"`
	Method   string `json:"method"`
	Detected bool   `json:"detected"`
}

// Task represents a scan task
type Task struct {
	ID            string       `json:"id"`
	FileName      string       `json:"fileName"`
	FilePath      string       `json:"filePath"`
	FileSize      int64        `json:"fileSize"`
	Status        string       `json:"status"`
	AnalysisID    string       `json:"analysisId"`
	MD5           string       `json:"md5"`
	SHA256        string       `json:"sha256"`
	SHA1          string       `json:"sha1"`
	ScanTime      string       `json:"scanTime"`
	DetectionRate int          `json:"detectionRate"`
	Results       []ScanResult `json:"results"`
	Stats         StatsInfo    `json:"stats"`
	FileInfo      FileInfo     `json:"fileInfo"`
}

// UploadFile uploads a file to VirusTotal for scanning
func (v *VTService) UploadFile(filePath string) (*Task, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// Create a buffer to write our multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://www.virustotal.com/api/v3/files", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("x-apikey", v.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var scanResp ScanResponse
	err = json.NewDecoder(resp.Body).Decode(&scanResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Create task
	task := &Task{
		ID:         scanResp.Data.ID,
		FileName:   filepath.Base(filePath),
		FilePath:   filePath,
		FileSize:   fileInfo.Size(),
		Status:     "queued",
		AnalysisID: scanResp.Data.ID,
		ScanTime:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return task, nil
}

// GetAnalysis gets the analysis result from VirusTotal
func (v *VTService) GetAnalysis(analysisID string) (*Task, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// Create the request
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/analyses/%s", analysisID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("x-apikey", v.apiKey)

	// Send the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var analysisResp AnalysisResponse
	err = json.NewDecoder(resp.Body).Decode(&analysisResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Create task with results
	task := &Task{
		ID:         analysisID,
		AnalysisID: analysisID,
		Status:     analysisResp.Data.Attributes.Status,
		MD5:        analysisResp.Meta.FileInfo.Md5,
		SHA1:       analysisResp.Meta.FileInfo.Sha1,
		SHA256:     analysisResp.Meta.FileInfo.Sha256,
		ScanTime:   time.Now().Format("2006-01-02 15:04:05"),
		Stats: StatsInfo{
			Malicious:        analysisResp.Data.Attributes.Stats.Malicious,
			Suspicious:       analysisResp.Data.Attributes.Stats.Suspicious,
			Harmless:         analysisResp.Data.Attributes.Stats.Harmless,
			Undetected:       analysisResp.Data.Attributes.Stats.Undetected,
			TypeUnsupported:  analysisResp.Data.Attributes.Stats.TypeUnsupported,
			ConfirmedTimeout: analysisResp.Data.Attributes.Stats.ConfirmedTimeout,
			Timeout:          analysisResp.Data.Attributes.Stats.Timeout,
			Failure:          analysisResp.Data.Attributes.Stats.Failure,
		},
		FileInfo: FileInfo{
			MD5:    analysisResp.Meta.FileInfo.Md5,
			SHA1:   analysisResp.Meta.FileInfo.Sha1,
			SHA256: analysisResp.Meta.FileInfo.Sha256,
			Size:   int64(analysisResp.Meta.FileInfo.Size),
		},
	}

	// Calculate detection rate
	total := task.Stats.Malicious + task.Stats.Suspicious + task.Stats.Harmless + task.Stats.Undetected
	if total > 0 {
		task.DetectionRate = int(float64(task.Stats.Malicious+task.Stats.Suspicious) / float64(total) * 100)
	}

	// Parse results
	task.Results = make([]ScanResult, 0)
	for engine, result := range analysisResp.Data.Attributes.Results {
		detected := result.Category == "malicious" || result.Category == "suspicious"
		task.Results = append(task.Results, ScanResult{
			Engine:   engine,
			Category: result.Category,
			Result:   result.Result,
			Method:   result.Method,
			Detected: detected,
		})
	}

	return task, nil
}

// OpenFileDialog opens a file dialog and returns the selected file path.
func (v *VTService) OpenFileDialog() string {
	result, err := v.app.Dialog.OpenFile().
		CanChooseFiles(true).
		SetTitle("请选择一个样本文件").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// OpenFileDialogs opens a directory dialog and returns the selected directory path.
func (v *VTService) OpenFileDialogs() string {
	result, err := v.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("请选择一个样本目录").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// GetFileInfo 获取文件信息
func (v *VTService) GetFileInfo(filePath string) (fileName string, fileSize int64, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get file info: %v", err)
	}
	return filepath.Base(filePath), fileInfo.Size(), nil
}

// VTTask VirusTotal 任务（用于前端展示）
type VTTask struct {
	ID            string       `json:"id"`
	TaskType      string       `json:"taskType"`      // 任务类型: single/directory
	FileName      string       `json:"fileName"`
	FilePath      string       `json:"filePath"`
	FileSize      int64        `json:"fileSize"`
	Status        string       `json:"status"`
	AnalysisID    string       `json:"analysisId"`
	CreatedAt     time.Time    `json:"createdAt"`
	CompletedAt   time.Time    `json:"completedAt"`
	MD5           string       `json:"md5"`
	SHA256        string       `json:"sha256"`
	SHA1          string       `json:"sha1"`
	DetectionRate int          `json:"detectionRate"`
	TotalEngines  int          `json:"totalEngines"`
	Results       []ScanResult `json:"results"`
	Stats         StatsInfo    `json:"stats"`
	FileInfo      FileInfo     `json:"fileInfo"`
	ScanTime      string       `json:"scanTime"`

	// 目录扫描专用字段
	TotalFiles      int `json:"totalFiles"`      // 目录下文件总数
	MaliciousFiles  int `json:"maliciousFiles"`  // 带毒文件数
	SuspiciousFiles int `json:"suspiciousFiles"` // 可疑文件数
	CompletedFiles  int `json:"completedFiles"`  // 已完成扫描文件数
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

// FileInfo 文件信息
type FileInfo struct {
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	Size   int64  `json:"size"`
}

// convertResults 转换检测结果格式
func convertResults(results []ScanResult) []store.ScanResult {
	if results == nil {
		return nil
	}
	converted := make([]store.ScanResult, len(results))
	for i, r := range results {
		converted[i] = store.ScanResult{
			Engine:   r.Engine,
			Category: r.Category,
			Result:   r.Result,
			Method:   r.Method,
		}
	}
	return converted
}

// CreateScanTask 创建扫描任务并上传文件
func (v *VTService) CreateScanTask(filePath string) (*VTTask, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// 生成任务ID
	taskID := generateTaskID()

	// 创建 VTTask
	vtTask := &VTTask{
		ID:         taskID,
		TaskType:   string(store.TaskTypeSingle),
		FileName:   filepath.Base(filePath),
		FilePath:   filePath,
		FileSize:   fileInfo.Size(),
		Status:     "pending",
		CreatedAt:  time.Now(),
		TotalFiles: 1,
	}

	// 创建对应的 VTFile 记录
	vtFile := &store.VTFile{
		ID:       generateFileID(),
		TaskID:   taskID,
		FileName: filepath.Base(filePath),
		FilePath: filePath,
		FileSize: fileInfo.Size(),
		Status:   "pending",
	}

	// 保存 VTFile
	if err := store.SaveVTFile(vtFile); err != nil {
		fmt.Printf("Failed to save VTFile: %v\n", err)
	}

	// 保存 VTTask
	tm := GetTaskManager()
	if err := tm.SaveTask(vtTask); err != nil {
		fmt.Printf("Failed to save task: %v\n", err)
	}

	// 上传文件到 VirusTotal
	task, err := v.UploadFile(filePath)
	if err != nil {
		vtTask.Status = "failed"
		vtFile.Status = "failed"
		tm.SaveTask(vtTask)
		store.SaveVTFile(vtFile)
		return nil, err
	}

	// 更新任务和文件状态
	vtTask.AnalysisID = task.AnalysisID
	vtTask.Status = "queued"
	vtTask.ScanTime = task.ScanTime

	vtFile.AnalysisID = task.AnalysisID
	vtFile.Status = "queued"

	tm.SaveTask(vtTask)
	store.SaveVTFile(vtFile)

	return vtTask, nil
}

// CheckAnalysisStatus 检查分析状态
func (v *VTService) CheckAnalysisStatus(taskID string) (*VTTask, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// 从数据库获取任务
	tm := GetTaskManager()
	vtTask, err := tm.LoadTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %v", err)
	}

	// 如果任务已完成，直接返回
	if vtTask.Status == "completed" {
		return vtTask, nil
	}

	// 调用 VirusTotal API 获取分析结果
	task, err := v.GetAnalysis(vtTask.AnalysisID)
	if err != nil {
		return nil, err
	}

	// 更新任务状态
	vtTask.Status = task.Status
	vtTask.MD5 = task.MD5
	vtTask.SHA256 = task.SHA256
	vtTask.SHA1 = task.SHA1
	vtTask.Stats = task.Stats
	vtTask.FileInfo = task.FileInfo
	vtTask.Results = task.Results
	vtTask.DetectionRate = task.DetectionRate
	vtTask.TotalEngines = task.Stats.Malicious + task.Stats.Suspicious + task.Stats.Harmless + 
		task.Stats.Undetected + task.Stats.TypeUnsupported

	// 如果分析完成，更新完成时间
	if task.Status == "completed" {
		vtTask.CompletedAt = time.Now()
		vtTask.CompletedFiles = 1
		vtTask.ScanTime = time.Now().Format("2006-01-02 15:04:05")

		// 更新对应的 VTFile 记录
		files, _ := store.GetVTFilesByTaskID(taskID)
		if len(files) > 0 {
			vtFile := files[0]
			vtFile.Status = "completed"
			vtFile.MD5 = task.MD5
			vtFile.SHA256 = task.SHA256
			vtFile.SHA1 = task.SHA1
			vtFile.Stats = store.StatsInfo(task.Stats)
			vtFile.Results = convertResults(task.Results)
			vtFile.DetectionRate = task.DetectionRate
			vtFile.TotalEngines = vtTask.TotalEngines
			vtFile.CompletedAt = time.Now()
			store.SaveVTFile(vtFile)
		}
	}

	// 保存到数据库
	if err := tm.SaveTask(vtTask); err != nil {
		fmt.Printf("Failed to update task: %v\n", err)
	}

	return vtTask, nil
}

// GetAllTasks 获取所有任务
func (v *VTService) GetAllTasks() ([]*VTTask, error) {
	tm := GetTaskManager()
	return tm.LoadAllTasks()
}

// DeleteTask 删除任务
func (v *VTService) DeleteTask(taskID string) error {
	tm := GetTaskManager()
	return tm.DeleteTask(taskID)
}

// GetTask 获取单个任务
func (v *VTService) GetTask(taskID string) (*VTTask, error) {
	tm := GetTaskManager()
	return tm.LoadTask(taskID)
}

// ========== 批量目录扫描功能 ==========

// startScanWorkers 启动扫描工作池
func (v *VTService) startScanWorkers() {
	for i := 0; i < v.concurrency; i++ {
		go v.scanWorker()
	}
}

// scanWorker 扫描工作协程
func (v *VTService) scanWorker() {
	for {
		select {
		case job := <-v.scanQueue:
			v.processScanJob(job)
			v.scanWg.Done()
		case <-v.scanStopChan:
			return
		}
	}
}

// processScanJob 处理扫描任务
func (v *VTService) processScanJob(job *scanJob) {
	if v.apiKey == "" {
		job.file.Status = "failed"
		store.SaveVTFile(job.file)
		return
	}

	// 上传文件
	task, err := v.UploadFile(job.file.FilePath)
	if err != nil {
		fmt.Printf("Failed to upload file %s: %v\n", job.file.FileName, err)
		job.file.Status = "failed"
		store.SaveVTFile(job.file)
		v.updateDirectoryTaskStats(job.taskID)
		return
	}

	// 更新文件状态
	job.file.AnalysisID = task.ID
	job.file.Status = "queued"
	store.SaveVTFile(job.file)

	// 轮询检查分析结果
	for i := 0; i < 60; i++ { // 最多等待10分钟
		time.Sleep(10 * time.Second)

		analysisTask, err := v.GetAnalysis(job.file.AnalysisID)
		if err != nil {
			fmt.Printf("Failed to get analysis for %s: %v\n", job.file.FileName, err)
			continue
		}

		if analysisTask.Status == "completed" {
			// 更新文件结果
			job.file.Status = "completed"
			job.file.MD5 = analysisTask.MD5
			job.file.SHA256 = analysisTask.SHA256
			job.file.SHA1 = analysisTask.SHA1
			job.file.DetectionRate = analysisTask.DetectionRate
			job.file.TotalEngines = analysisTask.Stats.Malicious + analysisTask.Stats.Suspicious + analysisTask.Stats.Harmless +
				analysisTask.Stats.Undetected + analysisTask.Stats.TypeUnsupported
			job.file.Stats = store.StatsInfo(analysisTask.Stats)
			job.file.Results = convertResults(analysisTask.Results)
			job.file.CompletedAt = time.Now()
			store.SaveVTFile(job.file)

			// 更新目录任务统计
			v.updateDirectoryTaskStats(job.taskID)
			return
		}
	}

	// 超时
	job.file.Status = "timeout"
	store.SaveVTFile(job.file)
	v.updateDirectoryTaskStats(job.taskID)
}

// updateDirectoryTaskStats 更新目录任务统计信息
func (v *VTService) updateDirectoryTaskStats(taskID string) {
	tm := GetTaskManager()
	vtTask, err := tm.LoadTask(taskID)
	if err != nil {
		return
	}

	if vtTask.TaskType != string(store.TaskTypeDirectory) {
		return
	}

	// 获取所有子文件
	files, err := store.GetVTFilesByTaskID(taskID)
	if err != nil {
		return
	}

	// 统计
	completedFiles := 0
	maliciousFiles := 0
	suspiciousFiles := 0

	for _, f := range files {
		if f.Status == "completed" {
			completedFiles++
			if f.Stats.Malicious > 0 {
				maliciousFiles++
			}
			if f.Stats.Suspicious > 0 {
				suspiciousFiles++
			}
		}
	}

	// 更新任务
	vtTask.CompletedFiles = completedFiles
	vtTask.MaliciousFiles = maliciousFiles
	vtTask.SuspiciousFiles = suspiciousFiles

	// 检查是否全部完成
	if completedFiles == vtTask.TotalFiles {
		vtTask.Status = "completed"
		vtTask.CompletedAt = time.Now()
	}

	tm.SaveTask(vtTask)
}

// CreateDirectoryScanTask 创建目录批量扫描任务
func (v *VTService) CreateDirectoryScanTask(directoryPath string) (*VTTask, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	// 收集目录下所有文件
	var files []string
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的文件
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in directory")
	}

	// 创建任务
	taskID := fmt.Sprintf("VT-DIR-%d", time.Now().UnixNano())
	vtTask := &VTTask{
		ID:            taskID,
		TaskType:      string(store.TaskTypeDirectory),
		FileName:      filepath.Base(directoryPath),
		FilePath:      directoryPath,
		Status:        "scanning",
		CreatedAt:     time.Now(),
		ScanTime:      time.Now().Format("2006-01-02 15:04:05"),
		TotalFiles:    len(files),
	}

	// 保存任务
	tm := GetTaskManager()
	if err := tm.SaveTask(vtTask); err != nil {
		return nil, err
	}

	// 创建子文件记录并加入扫描队列
	for _, filePath := range files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		vtFile := &store.VTFile{
			ID:        fmt.Sprintf("VT-FILE-%d-%d", time.Now().UnixNano(), len(filePath)),
			TaskID:    taskID,
			FileName:  filepath.Base(filePath),
			FilePath:  filePath,
			FileSize:  fileInfo.Size(),
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		if err := store.SaveVTFile(vtFile); err != nil {
			fmt.Printf("Failed to save file record: %v\n", err)
			continue
		}

		// 加入扫描队列
		v.scanWg.Add(1)
		v.scanQueue <- &scanJob{file: vtFile, taskID: taskID}
	}

	return vtTask, nil
}

// GetTaskFiles 获取任务下所有文件
func (v *VTService) GetTaskFiles(taskID string) ([]*store.VTFile, error) {
	return store.GetVTFilesByTaskID(taskID)
}

// GetFileDetail 获取文件详情
func (v *VTService) GetFileDetail(fileID string) (*store.VTFile, error) {
	file, err := store.GetVTFile(fileID)
	if err != nil {
		return nil, err
	}

	// 如果文件还在扫描中，尝试刷新状态
	if file.Status == "queued" || file.Status == "pending" {
		if file.AnalysisID != "" && v.apiKey != "" {
			analysisTask, err := v.GetAnalysis(file.AnalysisID)
			if err == nil && analysisTask.Status == "completed" {
				file.Status = "completed"
				file.MD5 = analysisTask.MD5
				file.SHA256 = analysisTask.SHA256
				file.SHA1 = analysisTask.SHA1
				file.DetectionRate = analysisTask.DetectionRate
				file.TotalEngines = analysisTask.Stats.Malicious + analysisTask.Stats.Suspicious + analysisTask.Stats.Harmless +
					analysisTask.Stats.Undetected + analysisTask.Stats.TypeUnsupported
				file.Stats = store.StatsInfo(analysisTask.Stats)
				file.Results = convertResults(analysisTask.Results)
				file.CompletedAt = time.Now()
				store.SaveVTFile(file)
			}
		}
	}

	return file, nil
}

// RefreshFileStatus 刷新文件扫描状态
func (v *VTService) RefreshFileStatus(fileID string) (*store.VTFile, error) {
	file, err := store.GetVTFile(fileID)
	if err != nil {
		return nil, err
	}

	if file.Status == "completed" {
		return file, nil
	}

	if file.AnalysisID == "" || v.apiKey == "" {
		return file, nil
	}

	analysisTask, err := v.GetAnalysis(file.AnalysisID)
	if err != nil {
		return nil, err
	}

	if analysisTask.Status == "completed" {
		file.Status = "completed"
		file.MD5 = analysisTask.MD5
		file.SHA256 = analysisTask.SHA256
		file.SHA1 = analysisTask.SHA1
		file.DetectionRate = analysisTask.DetectionRate
		file.TotalEngines = analysisTask.Stats.Malicious + analysisTask.Stats.Suspicious + analysisTask.Stats.Harmless +
			analysisTask.Stats.Undetected + analysisTask.Stats.TypeUnsupported
		file.Stats = store.StatsInfo(analysisTask.Stats)
		file.Results = convertResults(analysisTask.Results)
		file.CompletedAt = time.Now()
		store.SaveVTFile(file)

		// 更新目录任务统计
		v.updateDirectoryTaskStats(file.TaskID)
	} else {
		file.Status = analysisTask.Status
		store.SaveVTFile(file)
	}

	return file, nil
}
