package api

import (
	"encoding/json"
	"fmt"

	"github.com/kwaain/nakisama/lib/request"
)

// LoginInfo 是 GetLoginInfo 的返回结构体
type LoginInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

// GetLoginInfo 获取登录号信息, 返回 LoginInfo 结构体.
func GetLoginInfo() (LoginInfo, error) {
	response, err := request.Send("/get_login_info", "POST", nil)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("获取登录号信息失败: %v", err)
	}

	data, err := json.Marshal(response.Data)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("将 data 字段转换为 JSON 失败: %v", err)
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(data, &loginInfo)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("解析登录号信息失败: %v", err)
	}

	return loginInfo, nil
}
