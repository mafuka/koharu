package api

import (
	"github.com/kwaain/nakisama/lib/request"
)

type SendPrivateMsgParam struct {
	UserID     int64  `json:"user_id"`            // 用户 ID
	GroupID    int64  `json:"group_id,omitempty"` // 群组 ID（可选），当主动发起临时会话时代表来源群号。机器人必须是管理员或群主。
	Message    string `json:"message"`            // 消息内容
	AutoEscape bool   `json:"auto_escape"`        // 是否将消息内容作为纯文本发送，默认为 false。仅在 message 字段是字符串时有效。
}

type SendPrivateMsgReturn struct {
	MessageID float64 `json:"message_id"` // 送出消息的 ID
}

// SendPrivateMsg 向指定用户发送一条私聊消息。
func SendPrivateMsg(param SendPrivateMsgParam) (SendPrivateMsgReturn, error) {
	var res SendPrivateMsgReturn
	err := request.Do("/send_private_msg", "POST", param, &res)
	if err != nil {
		return SendPrivateMsgReturn{}, err
	}

	return res, nil
}
