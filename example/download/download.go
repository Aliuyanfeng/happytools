package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// 服务器地址
	ServerURL = "http://127.0.0.1:8888"
	// 下载接口路径
	DownloadPath = "/external/fileDownload"
)

// DownloadClient 下载客户端
type DownloadClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewDownloadClient 创建下载客户端
func NewDownloadClient(baseURL string) *DownloadClient {
	return &DownloadClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 0, // 不设置超时，支持大文件下载
		},
	}
}

// DownloadProgress 下载进度回调
type DownloadProgress func(downloaded, total int64, speed float64)

// DownloadFile 下载文件
// fileID: 文件ID
// outputPath: 输出目录路径（如果为空则使用当前目录）
// progress: 进度回调函数（可选）
// 返回值: 实际保存的文件路径, error
func (c *DownloadClient) DownloadFile(fileID, outputPath string, progress DownloadProgress) (string, error) {
	downloadURL := fmt.Sprintf("%s%s?file_id=%s", c.BaseURL, DownloadPath, fileID)

	// 创建请求
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return "", fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
	}

	// 从 Content-Disposition 获取文件名
	filename := parseFilenameFromContentDisposition(resp.Header.Get("Content-Disposition"))
	if filename == "" {
		filename = fmt.Sprintf("download_%s", fileID)
	}

	// 确定最终保存路径
	var finalPath string
	if outputPath == "" {
		finalPath = filename
	} else {
		// 检查 outputPath 是目录还是文件
		if info, err := os.Stat(outputPath); err == nil && info.IsDir() {
			// 是目录，使用解析出的文件名
			finalPath = fmt.Sprintf("%s%c%s", outputPath, os.PathSeparator, filename)
		} else {
			// 是文件路径，直接使用
			finalPath = outputPath
		}
	}

	// 获取文件总大小
	var totalSize int64
	contentRange := resp.Header.Get("Content-Range")
	if contentRange != "" {
		// 解析 Content-Range: bytes start-end/total
		parts := strings.Split(contentRange, "/")
		if len(parts) == 2 {
			totalSize, _ = strconv.ParseInt(parts[1], 10, 64)
		}
	} else {
		totalSize = resp.ContentLength
	}

	// 检查是否已有部分下载的文件（断点续传）
	var startByte int64 = 0
	if info, err := os.Stat(finalPath); err == nil {
		localSize := info.Size()

		// 如果本地文件大小等于服务器文件大小，说明已完整下载
		if localSize == totalSize && totalSize > 0 {
			// fmt.Printf("文件已存在且完整，跳过下载: %s\n", finalPath)
			// fmt.Printf("文件大小: %s\n", formatBytes(localSize))
			// return finalPath, nil
			startByte = 0
		} else {
			// 如果本地文件大小大于服务器文件大小，说明文件可能已损坏，重新下载
			if localSize > totalSize && totalSize > 0 {
				fmt.Printf("本地文件大小(%s)大于服务器文件大小(%s)，将重新下载\n",
					formatBytes(localSize), formatBytes(totalSize))
				startByte = 0
			} else if localSize > 0 {
				// 正常的断点续传
				startByte = localSize
				fmt.Printf("检测到已下载部分文件，从 %s 处继续下载...\n", formatBytes(startByte))
			}
		}

	}

	// 如果需要断点续传，重新发送带 Range 头的请求
	if startByte > 0 {
		req, err = http.NewRequest("GET", downloadURL, nil)
		if err != nil {
			return "", fmt.Errorf("创建断点续传请求失败: %w", err)
		}
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("断点续传请求失败: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			return "", fmt.Errorf("断点续传失败，服务器返回状态码: %d", resp.StatusCode)
		}
	}

	// 打开输出文件（追加模式）
	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 如果是新下载，跳转到文件开头
	if startByte == 0 {
		file.Seek(0, 0)
		file.Truncate(0)
	}

	fmt.Printf("开始下载: %s\n", downloadURL)
	fmt.Printf("文件名: %s\n", filename)
	fmt.Printf("文件大小: %s\n", formatBytes(totalSize))
	fmt.Println("----------------------------------------")

	// 带进度的复制
	startTime := time.Now()
	downloaded := startByte
	buf := make([]byte, 32*1024) // 32KB 缓冲区
	lastPrintTime := time.Now()
	var lastDownloaded int64 = startByte

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := file.Write(buf[:n])
			if writeErr != nil {
				return "", fmt.Errorf("写入文件失败: %w", writeErr)
			}
			downloaded += int64(n)

			// 每秒更新一次进度显示
			now := time.Now()
			if now.Sub(lastPrintTime) >= time.Second {
				elapsed := now.Sub(startTime).Seconds()
				var speed float64
				if elapsed > 0 {
					speed = float64(downloaded-lastDownloaded) / now.Sub(lastPrintTime).Seconds()
				}
				lastPrintTime = now
				lastDownloaded = downloaded

				// 打印进度
				printProgress(downloaded, totalSize, speed)

				// 回调
				if progress != nil {
					progress(downloaded, totalSize, speed)
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("读取响应失败: %w", err)
		}
	}

	// 最终进度
	elapsed := time.Since(startTime).Seconds()
	avgSpeed := float64(downloaded-startByte) / elapsed
	printProgress(downloaded, totalSize, avgSpeed)

	fmt.Println("\n----------------------------------------")
	fmt.Printf("下载完成! 总耗时: %s\n", time.Since(startTime).Round(time.Millisecond))
	fmt.Printf("保存位置: %s\n", finalPath)

	return finalPath, nil
}

// printProgress 打印进度条
func printProgress(downloaded, total int64, speed float64) {
	if total <= 0 {
		fmt.Printf("\r已下载: %s | 速度: %s/s", formatBytes(downloaded), formatBytes(int64(speed)))
		return
	}

	percent := float64(downloaded) / float64(total) * 100
	barWidth := 40
	filled := int(float64(barWidth) * float64(downloaded) / float64(total))
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("\r[%s] %.1f%% | %s/%s | %.1f MB/s",
		bar, percent, formatBytes(downloaded), formatBytes(total), speed/1024/1024)
}

// formatBytes 格式化字节大小
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// parseFilenameFromContentDisposition 从 Content-Disposition 响应头解析文件名
// 支持格式:
// - attachment; filename="文件名.zip"
// - attachment; filename*=UTF-8”%E6%96%87%E4%BB%B6%E5%90%8D.zip
// - attachment; filename="文件名.zip"; filename*=UTF-8”%E6%96%87%E4%BB%B6%E5%90%8D.zip
func parseFilenameFromContentDisposition(contentDisposition string) string {
	if contentDisposition == "" {
		return ""
	}

	// 分割各个参数
	parts := strings.Split(contentDisposition, ";")
	var filename string
	var filenameStar string

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 优先处理 filename* (RFC 5987 格式)
		if strings.HasPrefix(part, "filename*=") {
			value := strings.TrimPrefix(part, "filename*=")
			value = strings.TrimSpace(value)

			// 格式: UTF-8''%E6%96%87%E4%BB%B6%E5%90%8D.zip
			if strings.Contains(value, "''") {
				encParts := strings.SplitN(value, "''", 2)
				if len(encParts) == 2 {
					// 解码 URL 编码的文件名
					decoded, err := url.QueryUnescape(encParts[1])
					if err == nil {
						filenameStar = decoded
					}
				}
			}
		} else if strings.HasPrefix(part, "filename=") {
			value := strings.TrimPrefix(part, "filename=")
			value = strings.TrimSpace(value)

			// 去除引号
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}

			// 尝试解码 URL 编码（有些服务器会对 filename 也进行编码）
			if decoded, err := url.QueryUnescape(value); err == nil {
				filename = decoded
			} else {
				filename = value
			}
		}
	}

	// 优先返回 filename* (RFC 5987 格式支持 UTF-8)
	if filenameStar != "" {
		return filenameStar
	}

	return filename
}

func main() {
	// 创建客户端
	client := NewDownloadClient(ServerURL)

	// 文件ID
	fileID := "019cdab4bace7ac89128d426c8828146"

	// 输出目录（传空字符串表示使用当前目录，文件名从响应头获取）
	outputPath := ""

	// 下载文件
	savedPath, err := client.DownloadFile(fileID, outputPath, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n文件已保存到: %s\n", savedPath)
}
