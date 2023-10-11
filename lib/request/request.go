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
	Status  string      `json:"status"`
	RetCode int         `json:"retcode"`
	Msg     string      `json:"msg"`
	Wording string      `json:"wording"`
	Data    interface{} `json:"data"`
}

// Do 发送 HTTP 请求，并把响应中的 Data 解析到 result 结构体中。
func Do(url, method string, body interface{}, result interface{}) error {
	baseURL := conf.Get().Server.Post

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	fullURL := baseURL + url

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
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
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	response := &Response{
		Data: result,
	}
	err = json.Unmarshal(respBody, response)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	switch response.RetCode {
	case 0:
		// 调用成功
		return nil
	case 1:
		// 已提交 async 处理
		return nil
	default:
		// 操作失败
		errMsg := fmt.Sprintf("%s: %s - %s", url, response.Msg, response.Wording)
		return fmt.Errorf(errMsg)
	}
}

//
