package api

import "github.com/kwaain/nakisama/lib/request"

type SetFriendAddRequestParam struct {
	Flag    string `json:"flag"`    // 加好友请求的 flag（需从上报的数据中获得）
	Approve bool   `json:"approve"` // 是否同意请求
	Remark  string `json:"remark"`  // 添加后的好友备注（仅在同意时有效）
}

func SetFriendAddRequest(param SetFriendAddRequestParam) error {
	var res SendGroupMsgReturn
	err := request.Do("/set_friend_add_request", "POST", param, &res)
	if err != nil {
		return err
	}

	return nil
}
