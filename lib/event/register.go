package event

import (
	"reflect"
)

// eventMap maps Event types to their corresponding structures.
var eventMap = make(map[Type]reflect.Type)

// Register registers events to the eventMap so they can be parsed to.
// It actually uses reflection and Event interface.
//
// Note: Events should only be registered once at program startup.
//
//	 package main
//		func main() {event.Register()}
func Register() {
	events := []Event{
		/********* REGULAR MESSAGE **********/
		&FriendMsg{},
		&GroupMsg{},
		&TempMsg{},
		&StrangerMsg{},
		&OtherClientMsg{},

		/******* SYNCHRONIZED MESSAGE *******/
		&FriendSyncMsg{},
		&GroupSyncMsg{},
		&TempSyncMsg{},
		&StrangerSyncMsg{},

		/********** BOT SELF EVENT **********/
		&BotOnline{},
		&BotOfflineActive{},
		&BotOfflineForce{},
		&BotOfflineDrop{},
		&BotRelogin{},

		/*********** FRIEND EVENT ***********/
		&FriendInputStatusChanged{},
		&FriendNickChanged{},
		&FriendAdd{},
		&FriendDelete{},

		/*********** GROUP EVENT ************/
		&BotGroupPermissionChange{},
		&BotMute{},
		&BotUnmute{},
		&BotJoinGroup{},
		&BotLeaveKick{},
		&BotLeaveDisband{},
		&FriendRecall{},
		&Nudge{},
		&GroupNameChange{},
		&GroupEntAnnChange{},
		&GroupAllowAnonChat{},
		&GroupAllowConfessTalk{},
		&GroupAllowMemberInvite{},
		&MemberJoin{},
		&MemberLeaveQuit{},
		&MemberCardChange{},
		&MemberSpecialTitleChange{},
		&MemberPermissionChange{},
		&MemberMute{},
		&MemberUnmute{},
		&MemberHonorChange{},

		/********** Request Events **********/
		&NewFriendRequest{},
		&MemberJoinRequest{},
		&BotInvitedJoinGroupRequest{},

		/******** OTHER CLIENT EVENT ********/
		&OtherClientOnline{},
		&OtherClientOffline{},

		/********** COMMAND EVENT ***********/
		&CommandExecuted{},
	}

	for _, e := range events {
		t := e.GetType()
		eventMap[t] = reflect.TypeOf(e).Elem()
	}
}
