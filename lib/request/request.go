// request 包提供了发送 HTTP 请求的方法。
package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kwaain/nakisama/lib/conf"
)

// Response 是 API 响应的结构体。
type Response struct {
	Status  string                 `json:"status"`
	RetCode int                    `json:"retcode"`
	Msg     string                 `json:"msg"`
	Wording string                 `json:"wording"`
	Data    map[string]interface{} `json:"data"`
}

// Send 发送 HTTP 请求。
func Send(url, method string, body interface{}) (*Response, error) {
	baseURL := conf.Get().Server.Post

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	fullURL := baseURL + url

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = conf.Get().Server.Secret

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	response := &Response{}
	err = json.Unmarshal(respBody, response)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	switch response.RetCode {
	case 0:
		// 调用成功
		return response, nil
	case 1:
		// 已提交 async 处理
		return response, nil
	default:
		// 操作失败
		errMsg := fmt.Sprintf("%s: %s - %s", url, response.Msg, response.Wording)
		return nil, fmt.Errorf(errMsg)
	}
}
