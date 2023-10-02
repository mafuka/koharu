// msgpool 提供管理消息池的方法。
package msgpool

import (
	"fmt"

	"github.com/kwaain/nakisama/lib/db"
)

// PrivateMsg 表示私聊会话的消息表模型。
type PrivateMsg struct {
	db.PrivateMsg
}

// GroupMsg 表示群聊会话的消息表模型。
type GroupMsg struct {
	db.GroupMsg
}

// CheckIfPrivateChatExists 检查私聊会话是否存在。
func CheckIfPrivateChatExists(userID int64) bool {
	tableName := fmt.Sprintf("msg_p_%d", userID)
	return db.CheckIfTableExists(tableName)
}

// CheckIfGroupChatExists 检查群聊会话是否存在。
func CheckIfGroupChatExists(groupID int64) bool {
	tableName := fmt.Sprintf("msg_g_%d", groupID)
	return db.CheckIfTableExists(tableName)
}

// CreatePrivateChat 创建私聊会话。
func CreatePrivateChat(userID int64) (string, error) {
	tableName := fmt.Sprintf("msg_p_%d", userID)

	if CheckIfPrivateChatExists(userID) {
		return "", fmt.Errorf("private chat %d already exists", userID)
	}

	err := db.CreateTable(tableName, &PrivateMsg{})
	if err != nil {
		return "", err
	}

	return tableName, nil
}

// DropPrivateChat 删除私聊会话。
func DropPrivateChat(userID int64) error {
	tableName := fmt.Sprintf("msg_p_%d", userID)

	if !CheckIfPrivateChatExists(userID) {
		return fmt.Errorf("private chat %d does not exist", userID)
	}

	return db.DropTable(tableName)
}

// CreateGroupChat 创建群聊会话。
func CreateGroupChat(groupID int64) (string, error) {
	tableName := fmt.Sprintf("msg_g_%d", groupID)

	if CheckIfGroupChatExists(groupID) {
		return "", fmt.Errorf("group chat %d already exists", groupID)
	}

	err := db.CreateTable(tableName, &GroupMsg{})
	if err != nil {
		return "", err
	}

	return tableName, nil
}

// DropGroupChat 删除群聊会话。
func DropGroupChat(groupID int64) error {
	tableName := fmt.Sprintf("msg_g_%d", groupID)

	if !CheckIfGroupChatExists(groupID) {
		return fmt.Errorf("group chat %d does not exist", groupID)
	}

	return db.DropTable(tableName)
}
