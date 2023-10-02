package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	UserID    int64     `gorm:"primaryKey"` // 用户 ID
	Nickname  string    `gorm:"not null"`   // 用户昵称
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Group 群组表
type Group struct {
	gorm.Model
	GroupID   int64     `gorm:"primaryKey"` // 群组 ID
	GroupName string    `gorm:"not null"`   // 群组名
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// GroupMember 群成员表
type GroupMember struct {
	gorm.Model
	GroupID int64  `gorm:"primaryKey"` // 群组 ID
	UserID  int64  `gorm:"primaryKey"` // 用户 ID
	Card    string ``                  // 群名片
	Role    Role   ``                  // 成员角色
}

type Role int

const (
	RoleMember Role = iota // 0 群成员
	RoleAdmin              // 1 群管理员
	RoleOwner              // 2 群主
)

// Message 消息表
type Message struct {
	MessageID int64   `gorm:"primaryKey"` // 消息 ID
	Type      MsgType ``                  // 消息类型
	UserID    int64   `gorm:"index"`      // 用户 ID
	GroupID   int64   `gorm:"index"`      // 群组 ID
	Content   string  ``
}

// MessageType 消息类型
type MsgType int

const (
	SystemMsg  MsgType = iota // 系统消息
	PrivateMsg                // 私聊消息
	GroupMsg                  // 群聊消息
)
