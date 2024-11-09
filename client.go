package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"test-wss/common"
	"time"
)

func getServerAddr() string {
	server := os.Getenv("SERVER_ADDR")
	if server == "" {
		server = "localhost:80"
	}
	return server
}

func main() {
	conn, err := net.Dial("tcp", getServerAddr())
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	// 发送 WebSocket 升级请求
	fmt.Fprint(conn, "GET / HTTP/1.1\r\nHost: localhost\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==\r\nSec-WebSocket-Version: 13\r\n\r\n")

	// 使用 bufio.Reader 读取和处理响应
	reader := bufio.NewReader(conn)
	sessionID := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			common.SendFeishuMessage("Error reading: " + err.Error())
			return
		}
		if strings.HasPrefix(line, "HTTP/1.1 101") {
			continue // 跳过升级协议的头部信息
		}
		if strings.Contains(line, "WebSocket session ID") {
			sessionID = strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Println("Session ID received:", sessionID)
			break
		}
	}

	// 模拟每3秒发送一次消息
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := conn.Write([]byte("Ping from " + sessionID))
		if err != nil {
			fmt.Println("Error writing:", err)
			common.SendFeishuMessage("Session ID: " + sessionID + " Error writing: " + err.Error())
			return
		}
		fmt.Printf("Ping sent from session ID %s\n", sessionID)
	}
}
