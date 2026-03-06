// Package network
// @author: liuyanfeng
// @date: 2026/2/17
// @description: TCP/UDP协议调试服务
package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// TCPUDPConfig TCP/UDP配置
type TCPUDPConfig struct {
	Host    string `json:"host"`    // 主机地址
	Port    int    `json:"port"`    // 端口
	Timeout int    `json:"timeout"` // 超时时间(秒)
}

// MessageResult 消息结果
type MessageResult struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	Data    string `json:"data"`    // 数据
	HexData string `json:"hexData"` // 十六进制数据
	Length  int    `json:"length"`  // 数据长度
}

// ConnectionStatus 连接状态
type ConnectionStatus struct {
	IsConnected   bool   `json:"isConnected"`   // 是否已连接
	LocalAddr     string `json:"localAddr"`     // 本地地址
	RemoteAddr    string `json:"remoteAddr"`    // 远程地址
	Protocol      string `json:"protocol"`      // 协议类型
	Mode          string `json:"mode"`          // 模式: client/server
	ServerRunning bool   `json:"serverRunning"` // 服务端是否运行
}

// TCPUDPService TCP/UDP服务
type TCPUDPService struct {
	app            *application.App
	tcpConn        net.Conn
	udpClientConn  *net.UDPConn // UDP客户端连接
	udpServerConn  *net.UDPConn // UDP服务端连接
	mu             sync.Mutex
	isRunning      bool
	serverRunning  bool          // UDP服务端是否运行
	stopListenChan chan struct{} // 停止监听通道
}

// NewTCPUDPService 创建TCP/UDP服务实例
func NewTCPUDPService(app *application.App) *TCPUDPService {
	return &TCPUDPService{
		app:       app,
		isRunning: false,
	}
}

// TestTCPConnection 测试TCP连接
func (s *TCPUDPService) TestTCPConnection(config TCPUDPConfig) (*ConnectionResult, error) {
	startTime := time.Now()

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

// TestUDPConnection 测试UDP连接
func (s *TCPUDPService) TestUDPConnection(config TCPUDPConfig) (*ConnectionResult, error) {
	startTime := time.Now()

	// 设置默认超时
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	// 解析UDP地址
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return &ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("解析地址失败: %v", err),
			Latency: 0,
		}, nil
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
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

// ConnectTCP 建立TCP连接
func (s *TCPUDPService) ConnectTCP(config TCPUDPConfig) (*ConnectionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已连接,先关闭
	if s.tcpConn != nil {
		s.tcpConn.Close()
	}

	// 设置默认超时
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	// 建立TCP连接
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return &ConnectionStatus{
			IsConnected: false,
			Protocol:    "TCP",
		}, nil
	}

	s.tcpConn = conn
	s.isRunning = true

	return &ConnectionStatus{
		IsConnected: true,
		LocalAddr:   conn.LocalAddr().String(),
		RemoteAddr:  conn.RemoteAddr().String(),
		Protocol:    "TCP",
	}, nil
}

// DisconnectTCP 断开TCP连接
func (s *TCPUDPService) DisconnectTCP() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpConn != nil {
		err := s.tcpConn.Close()
		s.tcpConn = nil
		s.isRunning = false
		return err
	}

	return nil
}

// SendTCP 发送TCP数据
func (s *TCPUDPService) SendTCP(data string, isHex bool) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpConn == nil {
		return &MessageResult{
			Success: false,
			Message: "未建立TCP连接",
		}, nil
	}

	var sendData []byte
	var err error

	if isHex {
		// 十六进制字符串转字节数组
		sendData, err = hex.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return &MessageResult{
				Success: false,
				Message: fmt.Sprintf("十六进制数据格式错误: %v", err),
			}, nil
		}
	} else {
		sendData = []byte(data)
	}

	// 发送数据
	_, err = s.tcpConn.Write(sendData)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("发送失败: %v", err),
		}, nil
	}

	return &MessageResult{
		Success: true,
		Message: "发送成功",
		Data:    string(sendData),
		HexData: hex.EncodeToString(sendData),
		Length:  len(sendData),
	}, nil
}

// ReceiveTCP 接收TCP数据
func (s *TCPUDPService) ReceiveTCP(timeout int) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpConn == nil {
		return &MessageResult{
			Success: false,
			Message: "未建立TCP连接",
		}, nil
	}

	// 设置读取超时
	if timeout == 0 {
		timeout = 5
	}
	s.tcpConn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	// 读取数据
	buffer := make([]byte, 4096)
	n, err := s.tcpConn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return &MessageResult{
				Success: false,
				Message: "接收超时",
			}, nil
		}
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("接收失败: %v", err),
		}, nil
	}

	receivedData := buffer[:n]

	return &MessageResult{
		Success: true,
		Message: "接收成功",
		Data:    string(receivedData),
		HexData: hex.EncodeToString(receivedData),
		Length:  len(receivedData),
	}, nil
}

// SendUDP 发送UDP数据
func (s *TCPUDPService) SendUDP(config TCPUDPConfig, data string, isHex bool) (*MessageResult, error) {
	// 解析UDP地址
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("解析地址失败: %v", err),
		}, nil
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}
	defer conn.Close()

	var sendData []byte

	if isHex {
		// 十六进制字符串转字节数组
		sendData, err = hex.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return &MessageResult{
				Success: false,
				Message: fmt.Sprintf("十六进制数据格式错误: %v", err),
			}, nil
		}
	} else {
		sendData = []byte(data)
	}

	// 发送数据
	_, err = conn.Write(sendData)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("发送失败: %v", err),
		}, nil
	}

	return &MessageResult{
		Success: true,
		Message: "发送成功",
		Data:    string(sendData),
		HexData: hex.EncodeToString(sendData),
		Length:  len(sendData),
	}, nil
}

// ReceiveUDP 接收UDP数据
func (s *TCPUDPService) ReceiveUDP(config TCPUDPConfig, timeout int) (*MessageResult, error) {
	// 解析UDP地址
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("解析地址失败: %v", err),
		}, nil
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}
	defer conn.Close()

	// 设置读取超时
	if timeout == 0 {
		timeout = 5
	}
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	// 读取数据
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return &MessageResult{
				Success: false,
				Message: "接收超时",
			}, nil
		}
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("接收失败: %v", err),
		}, nil
	}

	receivedData := buffer[:n]

	return &MessageResult{
		Success: true,
		Message: "接收成功",
		Data:    string(receivedData),
		HexData: hex.EncodeToString(receivedData),
		Length:  len(receivedData),
	}, nil
}

// StringToHex 字符串转十六进制
func (s *TCPUDPService) StringToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// HexToString 十六进制转字符串
func (s *TCPUDPService) HexToString(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(strings.TrimSpace(hexStr))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetConnectionStatus 获取连接状态
func (s *TCPUDPService) GetConnectionStatus() *ConnectionStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpConn != nil {
		return &ConnectionStatus{
			IsConnected: true,
			LocalAddr:   s.tcpConn.LocalAddr().String(),
			RemoteAddr:  s.tcpConn.RemoteAddr().String(),
			Protocol:    "TCP",
			Mode:        "client",
		}
	}

	return &ConnectionStatus{
		IsConnected:   false,
		Protocol:      "TCP",
		ServerRunning: s.serverRunning,
	}
}

// StartUDPServer 启动UDP服务端监听
func (s *TCPUDPService) StartUDPServer(host string, port int) (*ConnectionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果服务端已运行,先关闭
	if s.udpServerConn != nil {
		s.udpServerConn.Close()
	}

	// 默认监听所有网卡
	if host == "" {
		host = "0.0.0.0"
	}

	// 监听UDP端口
	addr := fmt.Sprintf("%s:%d", host, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return &ConnectionStatus{
			IsConnected:   false,
			Protocol:      "UDP",
			Mode:          "server",
			ServerRunning: false,
		}, nil
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return &ConnectionStatus{
			IsConnected:   false,
			Protocol:      "UDP",
			Mode:          "server",
			ServerRunning: false,
		}, nil
	}

	s.udpServerConn = conn
	s.serverRunning = true
	s.stopListenChan = make(chan struct{})

	// 启动后台监听goroutine
	go s.listenUDPServer()

	return &ConnectionStatus{
		IsConnected:   true,
		LocalAddr:     conn.LocalAddr().String(),
		Protocol:      "UDP",
		Mode:          "server",
		ServerRunning: true,
	}, nil
}

// listenUDPServer 后台监听UDP数据
func (s *TCPUDPService) listenUDPServer() {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-s.stopListenChan:
			return
		default:
			// 设置读取超时，以便定期检查stop信号
			s.mu.Lock()
			if s.udpServerConn == nil {
				s.mu.Unlock()
				return
			}
			conn := s.udpServerConn
			s.mu.Unlock()

			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				// 连接已关闭或其他错误
				return
			}

			if n > 0 {
				receivedData := make([]byte, n)
				copy(receivedData, buffer[:n])

				// 通过事件推送到前端
				result := &MessageResult{
					Success: true,
					Message: fmt.Sprintf("接收成功 (来自: %s)", remoteAddr.String()),
					Data:    string(receivedData),
					HexData: hex.EncodeToString(receivedData),
					Length:  len(receivedData),
				}

				if s.app != nil {
					s.app.Event.Emit("network:udpReceived", result)
				}
			}
		}
	}
}

// StopUDPServer 停止UDP服务端
func (s *TCPUDPService) StopUDPServer() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 发送停止信号
	if s.stopListenChan != nil {
		close(s.stopListenChan)
		s.stopListenChan = nil
	}

	if s.udpServerConn != nil {
		err := s.udpServerConn.Close()
		s.udpServerConn = nil
		s.serverRunning = false
		return err
	}

	return nil
}

// SendUDPFromServer 从服务端发送UDP数据到指定客户端
func (s *TCPUDPService) SendUDPFromServer(targetHost string, targetPort int, data string, isHex bool) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.udpServerConn == nil {
		return &MessageResult{
			Success: false,
			Message: "UDP服务端未启动",
		}, nil
	}

	// 解析目标地址
	targetAddr := fmt.Sprintf("%s:%d", targetHost, targetPort)
	udpAddr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("解析目标地址失败: %v", err),
		}, nil
	}

	var sendData []byte

	if isHex {
		sendData, err = hex.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return &MessageResult{
				Success: false,
				Message: fmt.Sprintf("十六进制数据格式错误: %v", err),
			}, nil
		}
	} else {
		sendData = []byte(data)
	}

	// 发送数据
	_, err = s.udpServerConn.WriteToUDP(sendData, udpAddr)
	if err != nil {
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("发送失败: %v", err),
		}, nil
	}

	return &MessageResult{
		Success: true,
		Message: "发送成功",
		Data:    string(sendData),
		HexData: hex.EncodeToString(sendData),
		Length:  len(sendData),
	}, nil
}

// ReceiveUDPOnServer 在服务端接收UDP数据
func (s *TCPUDPService) ReceiveUDPOnServer(timeout int) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.udpServerConn == nil {
		return &MessageResult{
			Success: false,
			Message: "UDP服务端未启动",
		}, nil
	}

	// 设置读取超时
	if timeout == 0 {
		timeout = 5
	}
	s.udpServerConn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	// 读取数据
	buffer := make([]byte, 4096)
	n, remoteAddr, err := s.udpServerConn.ReadFromUDP(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return &MessageResult{
				Success: false,
				Message: "接收超时",
			}, nil
		}
		return &MessageResult{
			Success: false,
			Message: fmt.Sprintf("接收失败: %v", err),
		}, nil
	}

	receivedData := buffer[:n]

	return &MessageResult{
		Success: true,
		Message: fmt.Sprintf("接收成功 (来自: %s)", remoteAddr.String()),
		Data:    string(receivedData),
		HexData: hex.EncodeToString(receivedData),
		Length:  len(receivedData),
	}, nil
}

// GetUDPServerStatus 获取UDP服务端状态
func (s *TCPUDPService) GetUDPServerStatus() *ConnectionStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.udpServerConn != nil {
		return &ConnectionStatus{
			IsConnected:   true,
			LocalAddr:     s.udpServerConn.LocalAddr().String(),
			Protocol:      "UDP",
			Mode:          "server",
			ServerRunning: true,
		}
	}

	return &ConnectionStatus{
		IsConnected:   false,
		Protocol:      "UDP",
		Mode:          "server",
		ServerRunning: false,
	}
}
