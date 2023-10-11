package api

import (
	"github.com/kwaain/nakisama/lib/request"
)

type GetMsgParam struct {
	MessageID int32 `json:"message_id"`
}

type GetMsgReturn struct {
	Group       bool               `json:"group"`              // 是否是群消息
	GroupID     int64              `json:"group_id,omitempty"` // 是群消息时的群号
	MessageID   int32              `json:"message_id"`         // 消息 ID
	RealID      int32              `json:"real_id"`            // 消息真实 ID
	MessageType string             `json:"message_type"`       // 消息类型，group | private
	Sender      GetMsgReturnSender `json:"sender"`             // 发送者
	Time        int32              `json:"time"`               // 发送时间
	Message     string             `json:"message"`            // 消息内容
	RawMessage  string             `json:"raw_message"`        // 原始消息内容
}

// GetMsgReturnSender
type GetMsgReturnSender struct {
	Nickname string `json:"nickname"` // 发送者昵称
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
}

func GetMsg(param GetMsgParam) (GetMsgReturn, error) {
	var res GetMsgReturn
	err := request.Do("/get_msg", "POST", param, &res)
	if err != nil {
		return GetMsgReturn{}, err
	}

	return res, nil
}
