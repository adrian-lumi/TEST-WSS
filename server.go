package main

import (
	"fmt"
	"net/http"
	"test-wss/common"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// 定义 WebSocket 连接的结构体
type WebSocketConnection struct {
	Conn      *websocket.Conn
	Connected time.Time
}

// 存储所有活跃的 WebSocket 连接
var connections = make(map[string]*WebSocketConnection)

// WebSocket 升级器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有CORS请求
	},
}

// 处理 WebSocket 连接的函数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer wsConn.Close()

	// 生成唯一标识符，使用 UUID
	clientID := uuid.New().String()

	// 记录连接时间
	connectionInfo := &WebSocketConnection{
		Conn:      wsConn,
		Connected: time.Now(),
	}

	// 存储连接信息
	connections[clientID] = connectionInfo

	// 发送 clientID 给客户端
	if err := wsConn.WriteMessage(websocket.TextMessage, []byte(clientID)); err != nil {
		fmt.Println("Failed to send client ID:", err)
		common.SendFeishuMessage(fmt.Sprintf("Failed to send client ID: %s\n", err))
		return
	}

	fmt.Printf("New connection: %s at %v\n", clientID, connectionInfo.Connected)

	// 这里可以继续处理消息或其他逻辑
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			// 计算连接时长
			duration := time.Since(connectionInfo.Connected)
			fmt.Printf("[%s] 连接时长: %s, Error reading message: %s\n", clientID, duration, err)
			common.SendFeishuMessage(fmt.Sprintf("[%s] 连接时长: %s, Error reading message: %s\n", clientID, duration, err))
			break
		}
		fmt.Printf("Received message from %s: %s\n", clientID, string(message))
		response := "Received: " + string(message)
		wsConn.WriteMessage(websocket.TextMessage, []byte(response))
	}

	// 连接关闭时清理资源
	closeConnection(clientID)
}

// 清理连接的函数
func closeConnection(clientID string) {
	if conn, ok := connections[clientID]; ok {
		conn.Conn.Close()
		delete(connections, clientID)
		fmt.Printf("[%s] 回收 Connection closed\n", clientID)
	}
}

// 主函数，设置路由和启动服务器
func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/test-feishu-msg", func(w http.ResponseWriter, r *http.Request) {
		common.SendFeishuMessage("test-feishu-msg")
	})
	fmt.Println("Server started on :80")
	http.ListenAndServe(":80", nil)
}
