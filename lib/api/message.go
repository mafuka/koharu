package api

import (
	"errors"
	"fmt"

	"github.com/kwaain/nakisama/lib/request"
)

// PrivateMsg 表示私聊消息的结构体。
type PrivateMsg struct {
	UserID     int64  `json:"user_id"`            // 用户 ID
	GroupID    int64  `json:"group_id,omitempty"` // 群组 ID（可选），当主动发起临时会话时代表来源群号。机器人必须是管理员或群主。
	Message    string `json:"message"`            // 消息内容
	AutoEscape bool   `json:"auto_escape"`        // 是否将消息内容作为纯文本发送，默认为 false。仅在 message 字段是字符串时有效。
}

// SendPrivateMsg 向指定的 userID 发送一条私聊消息，内容为 message。
//
// U: 23-10-03 14:39
func SendPrivateMsg(msg PrivateMsg) (msgID float64, err error) {
	r, err := request.Send("/send_private_msg", "POST", msg)
	if err != nil {
		return 0, err
	}

	msgID, ok := r.Data["message_id"].(float64)
	if !ok {
		return 0, errors.New("无效的消息 ID 数据类型")
	}

	return msgID, nil
}

// GroupMsg 表示群聊消息的结构体。
type GroupMsg struct {
	GroupID    int64  `json:"group_id"`    // 群号
	Message    string `json:"message"`     // 消息内容
	AutoEscape bool   `json:"auto_escape"` // 是否将消息内容作为纯文本发送，默认为 false。仅在 message 字段是字符串时有效。
}

// SendGroupMsg 向指定的 groupID 发送一条私聊消息，内容为 message。
//
// U: 23-10-03 14:40
func SendGroupMsg(msg GroupMsg) (msgID float64, err error) {
	r, err := request.Send("/send_group_msg", "POST", msg)
	if err != nil {
		return 0, err
	}

	msgID, ok := r.Data["message_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("message_id 数据类型无效: %f", msgID)
	}

	return msgID, nil
}
