// Package network
// @author: liuyanfeng
// @date: 2026/3/5
// @description: DNS服务
package network

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// DNSService DNS服务
type DNSService struct {
	app *application.App
}

// NewDNSService 创建DNS服务
func NewDNSService() *DNSService {
	return &DNSService{}
}

// DNSFlushResult DNS刷新结果
type DNSFlushResult struct {
	Success     bool   `json:"success"`     // 是否成功
	Message     string `json:"message"`     // 消息
	Output      string `json:"output"`      // 命令输出
	AdminNeeded bool   `json:"adminNeeded"` // 是否需要管理员权限
}

// FlushDNS 刷新DNS缓存
func (d *DNSService) FlushDNS() (*DNSFlushResult, error) {
	var cmd *exec.Cmd
	var output bytes.Buffer

	switch runtime.GOOS {
	case "windows":
		// Windows: 使用 ipconfig /flushdns
		cmd = exec.Command("ipconfig", "/flushdns")
	case "darwin":
		// macOS: 使用 dscacheutil -flushcache 和 sudo killall -HUP mDNSResponder
		cmd = exec.Command("dscacheutil", "-flushcache")
	case "linux":
		// Linux: 尝试使用 systemd-resolve 或 resolvectl
		if _, err := exec.LookPath("resolvectl"); err == nil {
			cmd = exec.Command("resolvectl", "flush-caches")
		} else if _, err := exec.LookPath("systemd-resolve"); err == nil {
			cmd = exec.Command("systemd-resolve", "--flush-caches")
		} else {
			// 尝试使用 nscd
			cmd = exec.Command("nscd", "-i", "hosts")
		}
	default:
		return &DNSFlushResult{
			Success:     false,
			Message:     "不支持的操作系统",
			AdminNeeded: false,
		}, nil
	}

	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	outputStr := strings.TrimSpace(output.String())

	if err != nil {
		// 检查是否是权限问题
		errMsg := err.Error()
		adminNeeded := strings.Contains(errMsg, "access is denied") ||
			strings.Contains(errMsg, "permission denied") ||
			strings.Contains(errMsg, "requires elevation") ||
			strings.Contains(outputStr, "拒绝访问") ||
			strings.Contains(outputStr, "需要管理员权限")

		return &DNSFlushResult{
			Success:     false,
			Message:     "刷新DNS缓存失败: " + err.Error(),
			Output:      outputStr,
			AdminNeeded: adminNeeded,
		}, nil
	}

	return &DNSFlushResult{
		Success:     true,
		Message:     "DNS缓存已成功刷新",
		Output:      outputStr,
		AdminNeeded: false,
	}, nil
}

// GetDNSInfo 获取DNS信息
func (d *DNSService) GetDNSInfo() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if runtime.GOOS == "windows" {
		// Windows: 获取DNS配置
		cmd := exec.Command("ipconfig", "/all")
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output

		err := cmd.Run()
		if err != nil {
			result["error"] = err.Error()
			return result, nil
		}

		// 解析DNS服务器信息
		lines := strings.Split(output.String(), "\n")
		var dnsServers []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "DNS") && strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					dns := strings.TrimSpace(parts[1])
					if dns != "" && dns != "::" {
						dnsServers = append(dnsServers, dns)
					}
				}
			}
		}
		result["dnsServers"] = dnsServers
	}

	result["os"] = runtime.GOOS
	return result, nil
}

// SetApp 设置应用实例
func (d *DNSService) SetApp(app *application.App) {
	d.app = app
}
