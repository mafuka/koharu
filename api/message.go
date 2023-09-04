package api

import (
	"errors"

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
// 如果提供了 groupID 参数，则表示主动发起临时会话，groupID 为来源群号。
// 机器人必须是管理员或群主才能发起临时会话。
// autoEscape 参数控制消息内容是否作为纯文本发送，默认为 false。
//
// 如果调用成功，则返回消息的 ID。
func SendPrivateMsg(userID int64, groupID int64, message string, autoEscape bool) (float64, error) {
	params := &PrivateMsg{
		UserID:     userID,
		GroupID:    groupID,
		Message:    message,
		AutoEscape: autoEscape,
	}

	response, err := request.Send("/send_private_msg", "POST", params)
	if err != nil {
		return 0, err
	}

	messageID, ok := response.Data["message_id"].(float64)
	if !ok {
		return 0, errors.New("无效的消息 ID 数据类型")
	}

	return messageID, nil
}
