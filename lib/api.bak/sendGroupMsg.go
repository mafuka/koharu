package api

import (
	"github.com/kwaain/nakisama/lib/request"
)

type SendGroupMsgParam struct {
	GroupID    int64  `json:"group_id"`    // 群号
	Message    string `json:"message"`     // 消息内容
	AutoEscape bool   `json:"auto_escape"` // 是否将消息内容作为纯文本发送，默认为 false。仅在 message 字段是字符串时有效。
}

type SendGroupMsgReturn struct {
	MessageID float64 `json:"message_id"` // 消息 ID
}

// 发送群聊消息
func SendGroupMsg(param SendGroupMsgParam) (SendGroupMsgReturn, error) {
	var res SendGroupMsgReturn
	err := request.Do("/send_group_msg", "POST", param, &res)
	if err != nil {
		return SendGroupMsgReturn{}, err
	}

	return res, nil
}
