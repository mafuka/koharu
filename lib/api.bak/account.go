// api 包含了对 OneBot HTTP API 的封装。
package api

import (
	"github.com/kwaain/nakisama/lib/request"
)

type LoginInfoReturn struct {
	UserID   int64  `json:"user_id"`  // QQ 账号
	Nickname string `json:"nickname"` // 昵称
}

// GetLoginInfo 获取机器人的登录账号信息。
func GetLoginInfo() (LoginInfoReturn, error) {
	var res LoginInfoReturn
	err := request.Do("/get_login_info", "POST", nil, &res)
	if err != nil {
		return LoginInfoReturn{}, err
	}

	return res, nil
}
