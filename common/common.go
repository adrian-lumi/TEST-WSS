package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FeishuMessage 发送飞书消息的结构体
type FeishuMessage struct {
	MsgType string            `json:"msg_type"`
	Content map[string]string `json:"content"`
}

const (
	webhookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/1b9f7359-8f88-4bc4-ad46-d0aa41d5e5f9"
	token      = "ZmV6XamWqJTmIgEQUvGgSh"
)

// SendFeishuMessage 发送文本消息到飞书
func SendFeishuMessage(msg string) error {
	message := FeishuMessage{
		MsgType: "text", // 默认消息类型为文本
		Content: map[string]string{"text": msg},
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(msgBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // 设置 Authorization 头部

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	fmt.Println("Message sent successfully")

	return nil
}
