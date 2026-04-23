package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

// TCPClientInfo 已连接的 TCP 客户端信息
type TCPClientInfo struct {
	ID         string `json:"id"`
	RemoteAddr string `json:"remoteAddr"`
	ConnectedAt string `json:"connectedAt"`
}

// tcpServerClient 内部客户端管理
type tcpServerClient struct {
	id   string
	conn net.Conn
}

// StartTCPServer 启动 TCP 服务端
func (s *TCPUDPService) StartTCPServer(host string, port int) (*ConnectionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpListener != nil {
		return &ConnectionStatus{IsConnected: true, LocalAddr: s.tcpListenerAddr, Protocol: "TCP", Mode: "server", ServerRunning: true}, nil
	}

	if host == "" {
		host = "0.0.0.0"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return &ConnectionStatus{IsConnected: false, Protocol: "TCP", Mode: "server"}, fmt.Errorf("启动失败: %w", err)
	}

	s.tcpListener = ln
	s.tcpListenerAddr = ln.Addr().String()
	s.tcpServerClients = make(map[string]*tcpServerClient)
	s.tcpServerStop = make(chan struct{})

	go s.acceptTCPClients()

	return &ConnectionStatus{
		IsConnected:   true,
		LocalAddr:     s.tcpListenerAddr,
		Protocol:      "TCP",
		Mode:          "server",
		ServerRunning: true,
	}, nil
}

// StopTCPServer 停止 TCP 服务端
func (s *TCPUDPService) StopTCPServer() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpServerStop != nil {
		close(s.tcpServerStop)
		s.tcpServerStop = nil
	}
	if s.tcpListener != nil {
		s.tcpListener.Close()
		s.tcpListener = nil
	}
	for _, c := range s.tcpServerClients {
		c.conn.Close()
	}
	s.tcpServerClients = nil
	return nil
}

// GetTCPServerClients 获取已连接的客户端列表
func (s *TCPUDPService) GetTCPServerClients() []TCPClientInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []TCPClientInfo
	for _, c := range s.tcpServerClients {
		list = append(list, TCPClientInfo{ID: c.id, RemoteAddr: c.conn.RemoteAddr().String()})
	}
	return list
}

// SendTCPServerToClient 服务端向指定客户端发送数据（clientID 为空则广播）
func (s *TCPUDPService) SendTCPServerToClient(clientID string, data string, isHex bool) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sendData []byte
	var err error
	if isHex {
		sendData, err = hex.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return &MessageResult{Success: false, Message: "HEX 格式错误: " + err.Error()}, nil
		}
	} else {
		sendData = []byte(data)
	}

	sent := 0
	for id, c := range s.tcpServerClients {
		if clientID != "" && id != clientID {
			continue
		}
		if _, werr := c.conn.Write(sendData); werr == nil {
			sent++
		}
	}
	if sent == 0 {
		return &MessageResult{Success: false, Message: "无可用客户端"}, nil
	}
	return &MessageResult{
		Success: true, Message: fmt.Sprintf("已发送给 %d 个客户端", sent),
		Data: string(sendData), HexData: hex.EncodeToString(sendData), Length: len(sendData),
	}, nil
}

// acceptTCPClients 后台接受新连接
func (s *TCPUDPService) acceptTCPClients() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			select {
			case <-s.tcpServerStop:
				return
			default:
				return
			}
		}
		id := fmt.Sprintf("%s_%d", conn.RemoteAddr().String(), time.Now().UnixMilli())
		client := &tcpServerClient{id: id, conn: conn}

		s.mu.Lock()
		if s.tcpServerClients != nil {
			s.tcpServerClients[id] = client
		}
		s.mu.Unlock()

		if s.app != nil {
			s.app.Event.Emit("network:tcpClientConnected", TCPClientInfo{
				ID: id, RemoteAddr: conn.RemoteAddr().String(),
				ConnectedAt: time.Now().Format("15:04:05"),
			})
		}

		go s.handleTCPServerClient(client)
	}
}

// handleTCPServerClient 处理单个客户端的数据接收
func (s *TCPUDPService) handleTCPServerClient(client *tcpServerClient) {
	defer func() {
		client.conn.Close()
		s.mu.Lock()
		if s.tcpServerClients != nil {
			delete(s.tcpServerClients, client.id)
		}
		s.mu.Unlock()
		if s.app != nil {
			s.app.Event.Emit("network:tcpClientDisconnected", client.id)
		}
	}()

	buf := make([]byte, 4096)
	for {
		client.conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, err := client.conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 检查服务端是否已停止
				s.mu.Lock()
				stopped := s.tcpListener == nil
				s.mu.Unlock()
				if stopped {
					return
				}
				continue
			}
			return
		}
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			if s.app != nil {
				s.app.Event.Emit("network:tcpServerReceived", map[string]any{
					"clientID":   client.id,
					"remoteAddr": client.conn.RemoteAddr().String(),
					"data":       string(data),
					"hexData":    hex.EncodeToString(data),
					"length":     n,
					"time":       time.Now().Format("15:04:05"),
				})
			}
		}
	}
}

// StartTCPClient 建立 TCP 客户端连接并开始后台接收
func (s *TCPUDPService) StartTCPClient(host string, port int, timeout int) (*ConnectionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpConn != nil {
		s.tcpConn.Close()
		s.tcpConn = nil
	}
	if timeout == 0 {
		timeout = 10
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return &ConnectionStatus{IsConnected: false, Protocol: "TCP", Mode: "client"}, fmt.Errorf("连接失败: %w", err)
	}
	s.tcpConn = conn
	s.tcpClientStop = make(chan struct{})
	go s.receiveTCPClient()

	return &ConnectionStatus{
		IsConnected: true,
		LocalAddr:   conn.LocalAddr().String(),
		RemoteAddr:  conn.RemoteAddr().String(),
		Protocol:    "TCP",
		Mode:        "client",
	}, nil
}

// StopTCPClient 断开 TCP 客户端连接
func (s *TCPUDPService) StopTCPClient() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tcpClientStop != nil {
		close(s.tcpClientStop)
		s.tcpClientStop = nil
	}
	if s.tcpConn != nil {
		err := s.tcpConn.Close()
		s.tcpConn = nil
		return err
	}
	return nil
}

// SendTCPClient 客户端发送数据
func (s *TCPUDPService) SendTCPClient(data string, isHex bool) (*MessageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tcpConn == nil {
		return &MessageResult{Success: false, Message: "未连接"}, nil
	}
	var sendData []byte
	var err error
	if isHex {
		sendData, err = hex.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return &MessageResult{Success: false, Message: "HEX 格式错误: " + err.Error()}, nil
		}
	} else {
		sendData = []byte(data)
	}
	if _, err = s.tcpConn.Write(sendData); err != nil {
		return &MessageResult{Success: false, Message: "发送失败: " + err.Error()}, nil
	}
	return &MessageResult{
		Success: true, Message: "发送成功",
		Data: string(sendData), HexData: hex.EncodeToString(sendData), Length: len(sendData),
	}, nil
}

// receiveTCPClient 后台接收 TCP 客户端数据
func (s *TCPUDPService) receiveTCPClient() {
	buf := make([]byte, 4096)
	for {
		s.mu.Lock()
		conn := s.tcpConn
		s.mu.Unlock()
		if conn == nil {
			return
		}
		conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-s.tcpClientStop:
					return
				default:
					continue
				}
			}
			if s.app != nil {
				s.app.Event.Emit("network:tcpClientDisconnectedSelf", nil)
			}
			return
		}
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			if s.app != nil {
				s.app.Event.Emit("network:tcpClientReceived", map[string]any{
					"data":    string(data),
					"hexData": hex.EncodeToString(data),
					"length":  n,
					"time":    time.Now().Format("15:04:05"),
				})
			}
		}
	}
}

// 在 TCPUDPService 结构体中需要新增字段，通过扩展文件添加
