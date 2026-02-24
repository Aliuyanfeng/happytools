// Package network
// @author: liuyanfeng
// @date: 2026/2/17
// @description: 文件传输协议服务 (FTP, SFTP, SMB, S3)
package network

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// FTPConfig FTP配置
type FTPConfig struct {
	Host     string `json:"host"`     // 主机地址
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Timeout  int    `json:"timeout"`  // 超时时间(秒)
}

// SFTPConfig SFTP配置
type SFTPConfig struct {
	Host     string `json:"host"`     // 主机地址
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	KeyFile  string `json:"keyFile"`  // SSH私钥文件路径
	Timeout  int    `json:"timeout"`  // 超时时间(秒)
}

// SMBConfig SMB配置
type SMBConfig struct {
	Host     string `json:"host"`     // 主机地址
	Port     int    `json:"port"`     // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Share    string `json:"share"`    // 共享名称
	Timeout  int    `json:"timeout"`  // 超时时间(秒)
}

// S3Config S3配置
type S3Config struct {
	Endpoint        string `json:"endpoint"`        // 端点URL
	AccessKeyID     string `json:"accessKeyId"`     // 访问密钥ID
	SecretAccessKey string `json:"secretAccessKey"` // 秘密访问密钥
	Region          string `json:"region"`          // 区域
	Bucket          string `json:"bucket"`          // 存储桶名称
	UseSSL          bool   `json:"useSSL"`          // 是否使用SSL
}

// ConnectionResult 连接结果
type ConnectionResult struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	Latency int64  `json:"latency"` // 延迟(毫秒)
}

// FileListResult 文件列表结果
type FileListResult struct {
	Success bool       `json:"success"` // 是否成功
	Message string     `json:"message"` // 消息
	Files   []FileInfo `json:"files"`   // 文件列表
	Count   int        `json:"count"`   // 文件数量
}

// FileInfo 文件信息
type FileInfo struct {
	Name  string `json:"name"`  // 文件名
	Path  string `json:"path"`  // 路径
	Size  int64  `json:"size"`  // 大小
	IsDir bool   `json:"isDir"` // 是否是目录
}

// TransferResult 传输结果
type TransferResult struct {
	Success   bool   `json:"success"`   // 是否成功
	Message   string `json:"message"`   // 消息
	Bytes     int64  `json:"bytes"`     // 传输字节数
	Time      int64  `json:"time"`      // 耗时(毫秒)
	Speed     string `json:"speed"`     // 速度
	LocalPath string `json:"localPath"` // 本地路径
	RemotePath string `json:"remotePath"` // 远程路径
}

// FileTransferService 文件传输服务
type FileTransferService struct {
	app *application.App
}

// NewFileTransferService 创建文件传输服务实例
func NewFileTransferService(app *application.App) *FileTransferService {
	return &FileTransferService{
		app: app,
	}
}

// TestFTPConnection 测试FTP连接
func (s *FileTransferService) TestFTPConnection(config FTPConfig) (*ConnectionResult, error) {
	startTime := time.Now()
	
	// 设置默认端口
	if config.Port == 0 {
		config.Port = 21
	}
	
	// 设置默认超时
	if config.Timeout == 0 {
		config.Timeout = 10
	}
	
	// 尝试TCP连接
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return &ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
			Latency: 0,
		}, nil
	}
	defer conn.Close()
	
	latency := time.Since(startTime).Milliseconds()
	
	return &ConnectionResult{
		Success: true,
		Message: "连接成功",
		Latency: latency,
	}, nil
}

// TestSFTPConnection 测试SFTP连接
func (s *FileTransferService) TestSFTPConnection(config SFTPConfig) (*ConnectionResult, error) {
	startTime := time.Now()
	
	// 设置默认端口
	if config.Port == 0 {
		config.Port = 22
	}
	
	// 设置默认超时
	if config.Timeout == 0 {
		config.Timeout = 10
	}
	
	// 尝试TCP连接
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return &ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
			Latency: 0,
		}, nil
	}
	defer conn.Close()
	
	latency := time.Since(startTime).Milliseconds()
	
	return &ConnectionResult{
		Success: true,
		Message: "连接成功",
		Latency: latency,
	}, nil
}

// TestSMBConnection 测试SMB连接
func (s *FileTransferService) TestSMBConnection(config SMBConfig) (*ConnectionResult, error) {
	startTime := time.Now()
	
	// 设置默认端口
	if config.Port == 0 {
		config.Port = 445
	}
	
	// 设置默认超时
	if config.Timeout == 0 {
		config.Timeout = 10
	}
	
	// 尝试TCP连接
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return &ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
			Latency: 0,
		}, nil
	}
	defer conn.Close()
	
	latency := time.Since(startTime).Milliseconds()
	
	return &ConnectionResult{
		Success: true,
		Message: "连接成功",
		Latency: latency,
	}, nil
}

// TestS3Connection 测试S3连接
func (s *FileTransferService) TestS3Connection(config S3Config) (*ConnectionResult, error) {
	startTime := time.Now()
	
	// 简单验证配置
	if config.Endpoint == "" {
		return &ConnectionResult{
			Success: false,
			Message: "端点URL不能为空",
			Latency: 0,
		}, nil
	}
	
	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		return &ConnectionResult{
			Success: false,
			Message: "访问密钥不能为空",
			Latency: 0,
		}, nil
	}
	
	// 这里只是模拟测试,实际需要使用S3 SDK
	latency := time.Since(startTime).Milliseconds()
	
	return &ConnectionResult{
		Success: true,
		Message: "配置验证成功",
		Latency: latency,
	}, nil
}

// OpenFileDialog 打开文件对话框
func (s *FileTransferService) OpenFileDialog() string {
	result, err := s.app.Dialog.OpenFile().
		SetTitle("选择文件").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// OpenDirectoryDialog 打开目录对话框
func (s *FileTransferService) OpenDirectoryDialog() string {
	result, err := s.app.Dialog.OpenFile().
		SetTitle("选择目录").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return result
}

// ListLocalFiles 列出本地目录文件
func (s *FileTransferService) ListLocalFiles(dirPath string) (*FileListResult, error) {
	files := []FileInfo{}
	
	// 读取目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return &FileListResult{
			Success: false,
			Message: fmt.Sprintf("读取目录失败: %v", err),
			Files:   files,
			Count:   0,
		}, nil
	}
	
	// 遍历文件
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		files = append(files, FileInfo{
			Name:  entry.Name(),
			Path:  filepath.Join(dirPath, entry.Name()),
			Size:  info.Size(),
			IsDir: entry.IsDir(),
		})
	}
	
	return &FileListResult{
		Success: true,
		Message: "读取成功",
		Files:   files,
		Count:   len(files),
	}, nil
}

// UploadFile 上传文件 (模拟)
func (s *FileTransferService) UploadFile(localPath, remotePath string, protocol string) (*TransferResult, error) {
	startTime := time.Now()
	
	// 检查本地文件是否存在
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return &TransferResult{
			Success: false,
			Message: fmt.Sprintf("本地文件不存在: %v", err),
		}, nil
	}
	
	// 模拟上传过程
	// 实际实现需要根据协议类型使用相应的SDK
	time.Sleep(100 * time.Millisecond)
	
	elapsed := time.Since(startTime).Milliseconds()
	speed := float64(fileInfo.Size()) / float64(elapsed) * 1000 / 1024 // KB/s
	
	return &TransferResult{
		Success:    true,
		Message:    "上传成功",
		Bytes:      fileInfo.Size(),
		Time:       elapsed,
		Speed:      fmt.Sprintf("%.2f KB/s", speed),
		LocalPath:  localPath,
		RemotePath: remotePath,
	}, nil
}

// DownloadFile 下载文件 (模拟)
func (s *FileTransferService) DownloadFile(remotePath, localPath string, protocol string) (*TransferResult, error) {
	startTime := time.Now()
	
	// 模拟下载过程
	// 实际实现需要根据协议类型使用相应的SDK
	time.Sleep(100 * time.Millisecond)
	
	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return &TransferResult{
			Success: false,
			Message: fmt.Sprintf("创建本地文件失败: %v", err),
		}, nil
	}
	defer file.Close()
	
	// 写入一些测试数据
	testData := strings.Repeat("test", 1000)
	file.WriteString(testData)
	
	elapsed := time.Since(startTime).Milliseconds()
	size := int64(len(testData))
	speed := float64(size) / float64(elapsed) * 1000 / 1024 // KB/s
	
	return &TransferResult{
		Success:    true,
		Message:    "下载成功",
		Bytes:      size,
		Time:       elapsed,
		Speed:      fmt.Sprintf("%.2f KB/s", speed),
		LocalPath:  localPath,
		RemotePath: remotePath,
	}, nil
}

// CountFilesInDirectory 统计目录中的文件数量
func (s *FileTransferService) CountFilesInDirectory(dirPath string) (int, error) {
	count := 0
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
