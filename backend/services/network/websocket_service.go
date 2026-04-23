package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WSMessageResult WebSocket 消息结果
type WSMessageResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
	HexData string `json:"hexData"`
	Length  int    `json:"length"`
}

// WSConnectionStatus WebSocket 连接状态
type WSConnectionStatus struct {
	IsConnected bool   `json:"isConnected"`
	URL         string `json:"url"`
}

// WebSocketService WebSocket 客户端服务
type WebSocketService struct {
	app    *application.App
	mu     sync.Mutex
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	url    string
}

// NewWebSocketService 创建 WebSocket 服务实例
func NewWebSocketService(app *application.App) *WebSocketService {
	return &WebSocketService{app: app}
}

// Connect 连接到 WebSocket 服务器
func (s *WebSocketService) Connect(url string) (*WSConnectionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		s.cancel()
		s.conn.CloseNow()
		s.conn = nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, _, err := websocket.Dial(ctx, url, nil)
	cancel()
	if err != nil {
		return &WSConnectionStatus{IsConnected: false, URL: url}, fmt.Errorf("连接失败: %w", err)
	}

	recvCtx, recvCancel := context.WithCancel(context.Background())
	s.conn = conn
	s.ctx = recvCtx
	s.cancel = recvCancel
	s.url = url

	go s.receiveLoop(conn, recvCtx)

	return &WSConnectionStatus{IsConnected: true, URL: url}, nil
}

// Disconnect 断开连接
func (s *WebSocketService) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	if s.conn != nil {
		err := s.conn.Close(websocket.StatusNormalClosure, "bye")
		s.conn = nil
		return err
	}
	return nil
}

// SendText 发送文本消息
func (s *WebSocketService) SendText(data string) (*WSMessageResult, error) {
	s.mu.Lock()
	conn := s.conn
	ctx := s.ctx
	s.mu.Unlock()

	if conn == nil {
		return &WSMessageResult{Success: false, Message: "未连接"}, nil
	}
	if err := conn.Write(ctx, websocket.MessageText, []byte(data)); err != nil {
		return &WSMessageResult{Success: false, Message: "发送失败: " + err.Error()}, nil
	}
	return &WSMessageResult{
		Success: true, Message: "发送成功",
		Data: data, HexData: hex.EncodeToString([]byte(data)), Length: len(data),
	}, nil
}

// SendBinary 发送二进制（HEX 字符串）消息
func (s *WebSocketService) SendBinary(hexData string) (*WSMessageResult, error) {
	s.mu.Lock()
	conn := s.conn
	ctx := s.ctx
	s.mu.Unlock()

	if conn == nil {
		return &WSMessageResult{Success: false, Message: "未连接"}, nil
	}
	b, err := hex.DecodeString(hexData)
	if err != nil {
		return &WSMessageResult{Success: false, Message: "HEX 格式错误: " + err.Error()}, nil
	}
	if err := conn.Write(ctx, websocket.MessageBinary, b); err != nil {
		return &WSMessageResult{Success: false, Message: "发送失败: " + err.Error()}, nil
	}
	return &WSMessageResult{
		Success: true, Message: "发送成功",
		Data: string(b), HexData: hexData, Length: len(b),
	}, nil
}

// GetStatus 获取连接状态
func (s *WebSocketService) GetStatus() *WSConnectionStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &WSConnectionStatus{IsConnected: s.conn != nil, URL: s.url}
}

// receiveLoop 后台接收消息并推送事件
func (s *WebSocketService) receiveLoop(conn *websocket.Conn, ctx context.Context) {
	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			if s.app != nil {
				s.app.Event.Emit("network:wsDisconnected", map[string]any{"reason": err.Error()})
			}
			s.mu.Lock()
			if s.conn == conn {
				s.conn = nil
			}
			s.mu.Unlock()
			return
		}
		isText := msgType == websocket.MessageText
		payload := map[string]any{
			"data":    string(data),
			"hexData": hex.EncodeToString(data),
			"length":  len(data),
			"isText":  isText,
			"time":    time.Now().Format("15:04:05"),
		}
		if s.app != nil {
			s.app.Event.Emit("network:wsReceived", payload)
		}
	}
}
