package event

import (
	"reflect"
)

// eventMap maps Event types to their corresponding structures.
var eventMap = make(map[Type]reflect.Type)

// Register registers events to the eventMap so they can be parsed to.
// It actually uses reflection and structure tags.
//
// Note: Events should only be registered once at program startup.
//
//	 package main
//		func main() {event.Register()}
func Register() {
	types := []reflect.Type{
		/********* REGULAR MESSAGE **********/
		reflect.TypeOf(FriendMsg{}),
		reflect.TypeOf(GroupMsg{}),
		reflect.TypeOf(TempMsg{}),
		reflect.TypeOf(StrangerMsg{}),
		reflect.TypeOf(OtherClientMsg{}),

		/******* SYNCHRONIZED MESSAGE *******/
		reflect.TypeOf(FriendSyncMsg{}),
		reflect.TypeOf(GroupSyncMsg{}),
		reflect.TypeOf(TempSyncMsg{}),
		reflect.TypeOf(StrangerSyncMsg{}),

		/********** BOT SELF EVENT **********/
		reflect.TypeOf(BotOnline{}),
		reflect.TypeOf(BotOfflineActive{}),
		reflect.TypeOf(BotOfflineForce{}),
		reflect.TypeOf(BotOfflineDrop{}),
		reflect.TypeOf(BotRelogin{}),

		/*********** FRIEND EVENT ***********/
		reflect.TypeOf(FriendInputStatusChanged{}),
		reflect.TypeOf(FriendNickChanged{}),
		reflect.TypeOf(FriendAdd{}),
		reflect.TypeOf(FriendDelete{}),

		/*********** GROUP EVENT ************/
		reflect.TypeOf(BotGroupPermissionChange{}),
		reflect.TypeOf(BotMute{}),
		reflect.TypeOf(BotUnmute{}),
		reflect.TypeOf(BotJoinGroup{}),
		reflect.TypeOf(BotLeaveKick{}),
		reflect.TypeOf(BotLeaveDisband{}),
		reflect.TypeOf(FriendRecall{}),
		reflect.TypeOf(Nudge{}),
		reflect.TypeOf(GroupNameChange{}),
		reflect.TypeOf(GroupEntAnnChange{}),
		reflect.TypeOf(GroupAllowAnonChat{}),
		reflect.TypeOf(GroupAllowConfessTalk{}),
		reflect.TypeOf(GroupAllowMemberInvite{}),
		reflect.TypeOf(MemberJoin{}),
		reflect.TypeOf(MemberLeaveQuit{}),
		reflect.TypeOf(MemberCardChange{}),
		reflect.TypeOf(MemberSpecialTitleChange{}),
		reflect.TypeOf(MemberPermissionChange{}),
		reflect.TypeOf(MemberMute{}),
		reflect.TypeOf(MemberUnmute{}),
		reflect.TypeOf(MemberHonorChange{}),

		/********** Request Events **********/
		reflect.TypeOf(NewFriendRequest{}),
		reflect.TypeOf(MemberJoinRequest{}),
		reflect.TypeOf(BotInvitedJoinGroupRequest{}),

		/******** OTHER CLIENT EVENT ********/
		reflect.TypeOf(OtherClientOnline{}),
		reflect.TypeOf(OtherClientOffline{}),

		/********** COMMAND EVENT ***********/
		reflect.TypeOf(CommandExecuted{}),
	}

	for _, t := range types {
		eventType := t.Field(0).Tag.Get("type")
		eventMap[Type(eventType)] = t
	}
}
