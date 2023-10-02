// api 包含了对 OneBot HTTP API 的封装。
package api

import (
	"encoding/json"
	"fmt"

	"github.com/kwaain/nakisama/lib/request"
)

// LoginInfo 包含机器人的登录账号信息。
type LoginInfo struct {
	UserID   int64  `json:"user_id"`  // QQ 账号
	Nickname string `json:"nickname"` // 昵称
}

// GetLoginInfo 获取机器人的登录账号信息。
//
// 该接口无需参数。调用成功则返回 LoginInfo 结构体。
func GetLoginInfo() (LoginInfo, error) {
	response, err := request.Send("/get_login_info", "POST", nil)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("获取登录号信息失败: %v", err)
	}

	data, err := json.Marshal(response.Data)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("解析响应失败: %v", err)
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(data, &loginInfo)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("解析登录号信息失败: %v", err)
	}

	return loginInfo, nil
}
