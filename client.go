package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"test-wss/common"
	"time"

	"github.com/gorilla/websocket"
)

func getServerAddr() string {
	server := os.Getenv("SERVER_ADDR")
	if server == "" {
		server = "localhost:80" // 确保这里的端口是服务器监听的端口
	}
	return server
}

func getScheme() string {
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		scheme = "ws"
	}
	return scheme
}

func getHttpHost() string {
	return os.Getenv("HTTP_HOST")
}

func main() {
	serverAddr := getServerAddr()
	u := url.URL{Scheme: getScheme(), Host: serverAddr, Path: "/ws"}
	fmt.Printf("连接到 %s\n", u.String())

	mockEnv := os.Getenv("MOCK_ENV")

	// 创建一个自定义的Dialer，其中包含跳过证书验证的TLS配置
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}

	header := http.Header{}
	if getHttpHost() != "" {
		header.Set("Host", getHttpHost())
	}
	header.Set("ClientEnv", mockEnv)

	c, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		fmt.Printf("[%s]连接失败: %s\n", mockEnv, err)
		return
	}
	defer c.Close()

	// 接收从服务器发送的 UUID
	_, message, err := c.ReadMessage()
	if err != nil {
		fmt.Println("读取消息失败:", err)
		common.SendFeishuMessage(fmt.Sprintf("读取消息失败: %s\n", err))
		return
	}
	clientID := string(message) // 存储 clientID
	fmt.Printf("从服务器接收到的 client ID: %s\n", clientID)

	// 模拟每3秒发送一次消息
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Ping, %s", mockEnv)))
		if err != nil {
			fmt.Printf("[%s]发送消息时出错, client ID %s: %s\n", mockEnv, clientID, err)
			common.SendFeishuMessage(fmt.Sprintf("[%s]发送消息时出错, client ID %s: %s\n", mockEnv, clientID, err))
			return
		}
		fmt.Printf("[%s]发送 Ping, client ID %s\n", mockEnv, clientID)
	}
}
