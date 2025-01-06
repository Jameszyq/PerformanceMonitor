package utils

import (
	"PerformanceMonitor/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebHookInfo2 struct {
	Content string
}
type WebHookInfo struct {
	MsgType string
	Text    WebHookInfo2
}

func SendMsg(param WebHookInfo) (bool, string) {
	jsonPayload, err := json.Marshal(param)
	if err != nil {
		return false, fmt.Sprintf("json marshal error: %s", err.Error())
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", model.Config.WebHookAddress, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, fmt.Sprintf("Send Msg Error: %s", err.Error())
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Error sending request: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body", err.Error())
		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode == http.StatusOK {
		return true, ""
	} else {
		return false, fmt.Sprintf("Error sending request: %s", resp.Status)
	}
}
