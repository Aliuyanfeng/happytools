// Package services
// @author: liuyanfeng
// @date: 2025/8/4 17:24
// @description:
package monitor

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SysInfoService 提供主机资源信息
type SysInfoService struct {
	lastStats     map[string]net.IOCountersStat
	lastTimestamp time.Time
	statsMutex    sync.Mutex
}

// NewSysInfoService 创建新服务实例
func NewSysInfoService() *SysInfoService {
	return &SysInfoService{
		lastStats:     make(map[string]net.IOCountersStat),
		lastTimestamp: time.Time{},
	}
}

// HostInfo 主机信息
type HostInfo struct {
	Hostname       string `json:"hostname"`        // 主机名
	OS             string `json:"os"`              // 操作系统
	Platform       string `json:"platform"`        // 平台
	KernelVersion  string `json:"kernel_version"`  // 内核版本
	Architecture   string `json:"architecture"`    // 架构
	UptimeSeconds  uint64 `json:"uptime_seconds"`  // 运行时间（秒）
	UptimeReadable string `json:"uptime_readable"` // 运行时间（可读）
}

// CPUInfo CPU 使用率信息
type CPUInfo struct {
	CoreUsages []float64    `json:"core_usages"` // 每核使用率（单位：百分比）
	CoreInfo   cpu.InfoStat `json:"core_info"`   // CPU 信息
}

// MemoryInfo 内存使用信息
type MemoryInfo struct {
	Total       string  `json:"total"`        // 总内存（字节）
	Used        string  `json:"used"`         // 已使用（字节）
	UsedPercent float64 `json:"used_percent"` // 使用率（百分比）
}

// DiskInfo 磁盘使用信息
type DiskInfo struct {
	Path        string  `json:"path"`         // 挂载路径
	Total       uint64  `json:"total"`        // 总容量（字节）
	Used        uint64  `json:"used"`         // 已使用（字节）
	UsedPercent float64 `json:"used_percent"` // 使用率（百分比）
}

// NetworkInterface 网络接口信息
type NetworkInterface struct {
	Name            string  `json:"name"`            // 接口名称
	MAC             string  `json:"mac"`             // MAC 地址
	MTU             int     `json:"mtu"`             // 最大传输单元
	IsUp            bool    `json:"isUp"`            // 是否启用
	IPv4            string  `json:"ipv4"`            // IPv4 地址
	IPv6            string  `json:"ipv6"`            // IPv6 地址
	BytesSent       uint64  `json:"bytesSent"`       // 发送字节数
	BytesRecv       uint64  `json:"bytesRecv"`       // 接收字节数
	PacketsSent     uint64  `json:"packetsSent"`     // 发送包数
	PacketsRecv     uint64  `json:"packetsRecv"`     // 接收包数
	BytesSentRate   float64 `json:"bytesSentRate"`   // 每秒发送字节数
	BytesRecvRate   float64 `json:"bytesRecvRate"`   // 每秒接收字节数
	PacketsSentRate float64 `json:"packetsSentRate"` // 每秒发送包数
	PacketsRecvRate float64 `json:"packetsRecvRate"` // 每秒接收包数
}

type SysInfo struct {
	HostInfo          *HostInfo          `json:"host_info"`          // 主机信息
	CPUInfo           *CPUInfo           `json:"cpu_info"`           // CPU 使用率信息
	MemoryInfo        *MemoryInfo        `json:"memory_info"`        // 内存使用信息
	DiskInfo          *DiskInfo          `json:"disk_info"`          // 磁盘使用信息
	NetworkInterfaces []NetworkInterface `json:"network_interfaces"` // 网络接口信息
}

// GetSysInfo 定时获取系统信息
func (s *SysInfoService) GetSysInfo() (*SysInfo, error) {
	// 获取主机信息
	hostInfo, err := s.GetHostInfo()
	if err != nil {
		return nil, err
	}
	// 获取 CPU 信息
	cpuInfo, err := s.GetCPUInfo()
	if err != nil {
		return nil, err
	}
	memoryInfo, err := s.GetMemoryInfo()
	if err != nil {
		return nil, err
	}
	diskInfo, err := s.GetDiskInfo()
	if err != nil {
		return nil, err
	}
	networkInterfaces, err := s.GetNetworkInterfaces()
	if err != nil {
		return nil, err
	}
	return &SysInfo{
		HostInfo:          hostInfo,
		CPUInfo:           cpuInfo,
		MemoryInfo:        memoryInfo,
		DiskInfo:          diskInfo,
		NetworkInterfaces: networkInterfaces,
	}, nil
}

// GetHostInfo 获取主机信息
func (s *SysInfoService) GetHostInfo() (*HostInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &HostInfo{
		Hostname:       info.Hostname,
		OS:             info.OS,
		Platform:       info.Platform,
		KernelVersion:  info.KernelVersion,
		Architecture:   info.KernelArch,
		UptimeSeconds:  info.Uptime,
		UptimeReadable: time.Duration(info.Uptime * uint64(time.Second)).String(),
	}, nil
}

// GetCPUInfo 获取 CPU 使用率信息（按核心）
func (s *SysInfoService) GetCPUInfo() (*CPUInfo, error) {
	usages, err := cpu.Percent(time.Second, true) // true 表示分核心
	if err != nil {
		return nil, err
	}
	info, err := cpu.Info()
	if err != nil {
		fmt.Printf("CPU 获取失败: %v\n", err)
		return nil, err
	}
	roundedUsages := make([]float64, len(usages))
	for i, usage := range usages {
		roundedUsages[i], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
	}
	return &CPUInfo{
		CoreUsages: roundedUsages,
		CoreInfo:   info[0],
	}, nil
}

// GetMemoryInfo 获取内存信息
func (s *SysInfoService) GetMemoryInfo() (*MemoryInfo, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return &MemoryInfo{
		Total:       ByteSize(int64(vm.Total)),
		Used:        ByteSize(int64(vm.Used)),
		UsedPercent: vm.UsedPercent,
	}, nil
}

// GetDiskInfo 获取根目录磁盘使用情况
func (s *SysInfoService) GetDiskInfo() (*DiskInfo, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", usage.UsedPercent), 64)
	return &DiskInfo{
		Path:        usage.Path,
		Total:       usage.Total,
		Used:        usage.Used,
		UsedPercent: usedPercent,
	}, nil
}

// GetNetworkInterfaces 获取网络接口信息
func (s *SysInfoService) GetNetworkInterfaces() ([]NetworkInterface, error) {
	s.statsMutex.Lock()
	defer s.statsMutex.Unlock()

	currentTime := time.Now()
	var timeDelta float64

	if !s.lastTimestamp.IsZero() {
		timeDelta = currentTime.Sub(s.lastTimestamp).Seconds()
	} else {
		timeDelta = 0 // 第一次调用，无法计算速率
	}

	// 获取网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	// 获取网络统计信息
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network counters: %v", err)
	}

	// 构建结果
	var result []NetworkInterface
	for _, iface := range interfaces {
		// 查找对应的计数器
		var counter *net.IOCountersStat
		for _, c := range counters {
			if c.Name == iface.Name {
				counter = &c
				break
			}
		}

		// 提取 IP 地址
		var ipv4, ipv6 string
		for _, addr := range iface.Addrs {
			if addr.Addr != "" {
				if isIPv4(addr.Addr) {
					ipv4 = addr.Addr
				} else if isIPv6(addr.Addr) {
					ipv6 = addr.Addr
				}
			}
		}
		isUp := false
		for _, flag := range iface.Flags {
			if strings.Contains(flag, "up") {
				isUp = true
				break
			}
		}
		// 构建接口信息
		info := NetworkInterface{
			Name: iface.Name,
			MAC:  iface.HardwareAddr,
			MTU:  iface.MTU,
			IsUp: isUp,
			IPv4: ipv4,
			IPv6: ipv6,
		}

		// 添加计数器信息
		if counter != nil {
			info.BytesSent = counter.BytesSent
			info.BytesRecv = counter.BytesRecv
			info.PacketsSent = counter.PacketsSent
			info.PacketsRecv = counter.PacketsRecv

			// 计算速率
			if lastCounter, exists := s.lastStats[iface.Name]; exists && timeDelta > 0 {
				info.BytesSentRate = float64(counter.BytesSent-lastCounter.BytesSent) / timeDelta
				info.BytesRecvRate = float64(counter.BytesRecv-lastCounter.BytesRecv) / timeDelta
				info.PacketsSentRate = float64(counter.PacketsSent-lastCounter.PacketsSent) / timeDelta
				info.PacketsRecvRate = float64(counter.PacketsRecv-lastCounter.PacketsRecv) / timeDelta
			}
		}

		result = append(result, info)
	}

	// 更新缓存
	if timeDelta > 0 {
		for _, c := range counters {
			s.lastStats[c.Name] = c
		}
		s.lastTimestamp = currentTime
	}

	return result, nil
}

// 辅助函数：检查是否为 IPv4 地址
func isIPv4(addr string) bool {
	return len(addr) > 0 && addr[0] != ':' && addr[0] != '['
}

// 辅助函数：检查是否为 IPv6 地址
func isIPv6(addr string) bool {
	return len(addr) > 0 && (addr[0] == ':' || addr[0] == '[')
}

// ByteSize 将字节转换为人类可读的字符串表示
func ByteSize(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	unitIndex := 0
	size := float64(bytes)

	// 不断除以1024直到找到合适的单位
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	// 格式化数字，保留1位小数
	str := fmt.Sprintf("%.1f", size)

	// 去除不必要的零和小数点
	if strings.Contains(str, ".") {
		str = strings.TrimSuffix(str, ".0")
		str = strings.TrimSuffix(str, ".")
	}

	return str + units[unitIndex]
}
