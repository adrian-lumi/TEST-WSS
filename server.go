package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"test-wss/common"

	"github.com/google/uuid"
)

func handleConnection(conn net.Conn) {
	sessionID := uuid.New().String() // 为每个连接生成一个唯一的会话ID
	defer conn.Close()
	fmt.Printf("[%s] Connection accepted\n", sessionID)

	// 读取 HTTP 请求
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Printf("[%s] Error reading request: %v\n", sessionID, err)
		common.SendFeishuMessage("Error reading request: " + err.Error())
		return
	}

	// 检查是否为 WebSocket 升级请求
	if strings.ToLower(request.Header.Get("Upgrade")) == "websocket" {
		// 这里简化处理，实际应用中需要验证更多的头部信息
		fmt.Fprintf(conn, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n")
		fmt.Printf("[%s] WebSocket upgrade completed\n", sessionID)

		fmt.Fprintf(conn, "[%s] WebSocket session ID: %s\n", sessionID, sessionID)
		// 进入数据帧处理循环
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("[%s] Connection closed by client\n", sessionID)
				} else {
					fmt.Printf("[%s] Error reading from connection: %v\n", sessionID, err)
				}
				break
			}
			fmt.Printf("[%s] Received: %s\n", sessionID, string(buf[:n]))
		}
	} else {
		fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request")
		fmt.Printf("[%s] Bad WebSocket upgrade request\n", sessionID)
		common.SendFeishuMessage("Bad WebSocket upgrade request")
	}
}

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on :80")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}
