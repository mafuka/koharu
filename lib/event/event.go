package event

import (
	"encoding/json"

	"github.com/kwaain/nakisama/lib/msgchain"
)

// Type represents the event type.
//
// From now on, Type exists only to ensure type safety.
//
// // The event parser creates a mapping to the event structure for this constant.
type Type string

const (
	/********* REGULAR MESSAGE **********/

	// Friend Message
	TypeFriendMsg Type = "FriendMessage"

	// Group Message
	TypeGroupMsg Type = "GroupMessage"

	// Temporary Message in a Group
	TypeTempMsg Type = "TempMessage"

	// Stranger Message
	TypeStrangerMsg Type = "StrangerMessage"

	// Message from Other Clients
	TypeOtherClientMsg Type = "OtherClientMessage"

	/******* SYNCHRONIZED MESSAGE *******/

	// Synchronized Friend Message
	TypeFriendSyncMsg Type = "FriendSyncMessage"

	// Synchronized Group Message
	TypeGroupSyncMsg Type = "GroupSyncMessage"

	// Synchronized Temporary Message in a Group
	TypeTempSyncMsg Type = "TempSyncMessage"

	// Synchronized Stranger Message
	TypeStrangerSyncMsg Type = "StrangerSyncMessage"

	/********** BOT SELF EVENT **********/

	// Bot Online Event
	TypeBotOnlineEvent Type = "BotOnlineEvent"

	// Bot Active Offline Event
	TypeBotOfflineEventActive Type = "BotOfflineEventActive"

	// Bot Forced Offline Event
	TypeBotOfflineEventForce Type = "BotOfflineEventForce"

	// Bot Dropped Offline Event (due to server disconnection or network issues)
	TypeBotOfflineEventDropped Type = "BotOfflineEventDropped"

	// Bot Relogin Event
	TypeBotReloginEvent Type = "BotReloginEvent"

	/*********** FRIEND EVENT ***********/

	// Friend Input Status Changed Event
	TypeFriendInputStatusChangedEvent Type = "FriendInputStatusChangedEvent"

	// Friend Nickname Changed Event
	TypeFriendNickChangedEvent Type = "FriendNickChangedEvent"

	// Friend Added Event
	TypeFriendAddEvent Type = "FriendAddEvent"

	// Friend Deleted Event
	TypeFriendDeleteEvent Type = "FriendDeleteEvent"

	/*********** GROUP EVENT ************/

	// Bot Group Permission Change Event (the operator must be the group owner)
	TypeBotGroupPermissionChangeEvent Type = "BotGroupPermissionChangeEvent"

	// Bot Muted Event
	TypeBotMuteEvent Type = "BotMuteEvent"

	// Bot Unmuted Event
	TypeBotUnmuteEvent Type = "BotUnmuteEvent"

	// Bot Joined a New Group
	TypeBotJoinGroupEvent Type = "BotJoinGroupEvent"

	// Bot Active Leave Event
	TypeBotLeaveEventActive Type = "BotLeaveEventActive"

	// Bot Kicked Out of a Group
	TypeBotLeaveEventKick Type = "BotLeaveEventKick"

	// Bot Leave Event due to Group Dissolution (the operator must be the group owner)
	TypeBotLeaveEventDisband Type = "BotLeaveEventDisband"

	// Group Message Recall Event
	TypeGroupRecallEvent Type = "GroupRecallEvent"

	// Friend Message Recall Event
	TypeFriendRecallEvent Type = "FriendRecallEvent"

	// Nudge Event
	TypeNudgeEvent Type = "NudgeEvent"

	// Group Name Change Event
	TypeGroupNameChangeEvent Type = "GroupNameChangeEvent"

	// Group Entrance Announcement Change Event
	TypeGroupEntranceAnnouncementChangeEvent Type = "GroupEntranceAnnouncementChangeEvent"

	// Group Mute All Event
	TypeGroupMuteAllEvent Type = "GroupMuteAllEvent"

	// Group Anonymous Chat Allowed Event
	TypeGroupAllowAnonymousChatEvent Type = "GroupAllowAnonymousChatEvent"

	// Group Confess Talk Allowed Event
	TypeGroupAllowConfessTalkEvent Type = "GroupAllowConfessTalkEvent"

	// Group Member Invite Allowed Event
	TypeGroupAllowMemberInviteEvent Type = "GroupAllowMemberInviteEvent"

	// New Member Join Event
	TypeMemberJoinEvent Type = "MemberJoinEvent"

	// Member Kicked Out of the Group (the member is not the bot)
	TypeMemberLeaveEventKick Type = "MemberLeaveEventKick"

	// Member Active Leave Event (the member is not the bot)
	TypeMemberLeaveEventQuit Type = "MemberLeaveEventQuit"

	// Member Card Change Event
	TypeMemberCardChangeEvent Type = "MemberCardChangeEvent"

	// Member Special Title Change Event (only the group owner has the privilege to operate)
	TypeMemberSpecialTitleChangeEvent Type = "MemberSpecialTitleChangeEvent"

	// Member Permission Change Event (the member is not the bot)
	TypeMemberPermissionChangeEvent Type = "MemberPermissionChangeEvent"

	// Member Muted Event (the member is not the bot)
	TypeMemberMuteEvent Type = "MemberMuteEvent"

	// Member Unmuted Event (the member is not the bot)
	TypeMemberUnmuteEvent Type = "MemberUnmuteEvent"

	// Member Honor Change Event
	TypeMemberHonorChangeEvent Type = "MemberHonorChangeEvent"

	/********** Request Events **********/

	// New Friend Request Event
	TypeNewFriendRequestEvent Type = "NewFriendRequestEvent"

	// Member Join Request Event (Bot requires administrator permission)
	TypeMemberJoinRequestEvent Type = "MemberJoinRequestEvent"

	// Bot Invited Join Group Request Event
	TypeBotInvitedJoinGroupRequestEvent Type = "BotInvitedJoinGroupRequestEvent"

	/******** OTHER CLIENT EVENT ********/

	// Other Client Online Event
	TypeOtherClientOnlineEvent Type = "OtherClientOnlineEvent"

	// Other Client Offline Event
	TypeOtherClientOfflineEvent Type = "OtherClientOfflineEvent"

	/********** COMMAND EVENT ***********/

	// Command Executed Event
	TypeCommandExecutedEvent Type = "CommandExecutedEvent"
)

// Event has only one Type field for type matching.
// The final parsed event should be in the below type-specific structure .
type Event interface {
	GetType() Type
}

type BaseEvent struct {
	Type Type `json:"type"`
}

type MsgEvent interface {
	GetRawMsgChain() json.RawMessage
	SetMsgChain(msgchain.MsgChain)
}

/*********** COMMON TYPE ************/
/* Common types are offen reused in event types. */

// Friend represents a friend object, which is used in friend messages and friend events.
type Friend struct {
	ID       int    `json:"ID"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

// GroupMember typically represents a member of a group and is used in group-related messages.
type GroupMember struct {
	ID int `json:"id"`

	// MemberName, or "in-group nickname", which is the nickname used within the group
	MemberName string `json:"memberName"`

	// SpecialTitle is just "special title", within the group
	SpecialTitle string `json:"specialTitle"`

	// The member's permissions within the group
	Permission string `json:"permission"`

	// Timestamp when the member joined the group
	JoinTimestamp int `json:"joinTimestamp"`

	// Timestamp of the member's last message
	LastSpeakTimestamp int `json:"lastSpeakTimestamp"`

	// The remaining mute time when a member is muted
	MuteTimeRemaining int `json:"muteTimeRemaining"`

	// Group in which the member is a part of
	Group Group `json:"group"`
}

// GroupMemberSimple is another GroupMember with fewer fields.
type GroupMemberSimple struct {
	ID         int    `json:"id"`
	MemberName string `json:"memberName"`
	Permission string `json:"permission"`
	Group      Group  `json:"group"`
}

// Group is a generic structure often used to describe event-related group information.
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// Bot's permission in the group, OWNER, ADMINISTRATOR or MEMBER.
	Permission GroupPermission `json:"permission"`
}

// GroupPermission is a member's permission in the group,
// which can be OWNER, ADMINISTRATOR or MEMBER.
type GroupPermission string

const (
	// OWNER
	Owner GroupPermission = "OWNER"

	// ADMINISTRATOR
	Admin GroupPermission = "ADMINISTRATOR"

	// MEMBER
	Member GroupPermission = "MEMBER"
)

// Stranger represents a stranger and is used in stranger-related events.
type Stranger struct {
	ID       int    `json:"ID"`
	Nickname string `json:"nickname"`

	// (? idk why strangers would have remarks)
	Remark string `json:"remark"`
}

// OtherClient represents other clients of the same account and is used in OtherClientMsg.
type OtherClient struct {
	ID       int    `json:"id"`
	Platform string `json:"platform"`
}

/********* REGULAR MESSAGE **********/

// FriendMsg represents a message from a friend.
type FriendMsg struct {
	Type Type `json:"type" type:"FriendMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	Sender Friend `json:"sender"`
}

func (f *FriendMsg) GetType() Type {
	return TypeFriendMsg
}

func (f *FriendMsg) GetRawMsgChain() json.RawMessage {
	return f.RawMsgChain
}

func (f *FriendMsg) SetMsgChain(mc msgchain.MsgChain) {
	f.MsgChain = mc
}

// GroupMsg represents a message from a group.
type GroupMsg struct {
	Type Type `json:"type" type:"GroupMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	Sender GroupMember `json:"sender"`
}

func (g *Group) GetType() Type {
	return TypeGroupMsg
}

func (g *GroupMsg) GetRawMsgChain() json.RawMessage {
	return g.RawMsgChain
}

func (g *GroupMsg) SetMsgChain(mc msgchain.MsgChain) {
	g.MsgChain = mc
}

// TempMsg represents a temporary message in a group.
type TempMsg struct {
	Type Type `json:"type" type:"TempMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	Sender GroupMember `json:"sender"`
}

func (t *TempMsg) GetType() Type {
	return TypeTempMsg
}

func (t *TempMsg) GetRawMsgChain() json.RawMessage {
	return t.RawMsgChain
}

func (t *TempMsg) SetMsgChain(mc msgchain.MsgChain) {
	t.MsgChain = mc
}

// StrangerMsg represents a message from a stranger.
type StrangerMsg struct {
	Type Type `json:"type" type:"StrangerMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	Sender Stranger `json:"sender"`
}

func (s *StrangerMsg) GetType() Type {
	return TypeStrangerMsg
}

func (s *StrangerMsg) GetRawMsgChain() json.RawMessage {
	return s.RawMsgChain
}

func (s *StrangerMsg) SetMsgChain(mc msgchain.MsgChain) {
	s.MsgChain = mc
}

// OtherClientMsg represents a message from another client.
type OtherClientMsg struct {
	Type Type `json:"type" type:"OtherClientMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	Sender OtherClient `json:"sender"`
}

func (o *OtherClientMsg) GetType() Type {
	return TypeOtherClientMsg
}

func (o *OtherClientMsg) GetRawMsgChain() json.RawMessage {
	return o.RawMsgChain
}

func (o *OtherClientMsg) SetMsgChain(mc msgchain.MsgChain) {
	o.MsgChain = mc
}

/******* SYNCHRONIZED MESSAGE *******/
/* Synchronized messages are like normal messages, but are events generated when */
/* messages sent by other clients of the bot account are synchronized to the mirai. */
/* The sender of such events is always the bot itself, so they are omitted. */

// FriendSyncMsg is synchronized FriendMsg.
type FriendSyncMsg struct {
	Type Type `json:"type" type:"FriendSyncMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	// Subject is the target friend for sending.
	Subject Friend `json:"subject"`
}

func (f *FriendSyncMsg) GetType() Type {
	return TypeFriendSyncMsg
}

func (f *FriendSyncMsg) GetRawMsgChain() json.RawMessage {
	return f.RawMsgChain
}

func (f *FriendSyncMsg) SetMsgChain(mc msgchain.MsgChain) {
	f.MsgChain = mc
}

// GroupSyncMsg is synchronized GroupMsg.
type GroupSyncMsg struct {
	Type Type `json:"type" type:"GroupSyncMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	// Subject is the target group for sending.
	Subject Group `json:"subject"`
}

func (g *GroupSyncMsg) GetType() Type {
	return TypeGroupSyncMsg
}

func (g *GroupSyncMsg) GetRawMsgChain() json.RawMessage {
	return g.RawMsgChain
}

func (g *GroupSyncMsg) SetMsgChain(mc msgchain.MsgChain) {
	g.MsgChain = mc
}

// TempSyncMsg is synchronized TempMsg.
type TempSyncMsg struct {
	Type Type `json:"type" type:"TempSyncMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`

	// Subject is the target group member,
	// the corresponding group information is in the Group field of the group member.
	Subject GroupMember `json:"subject"`

	MsgChain msgchain.MsgChain
}

func (t *TempSyncMsg) GetType() Type {
	return TypeTempSyncMsg
}

func (t *TempSyncMsg) GetRawMsgChain() json.RawMessage {
	return t.RawMsgChain
}

func (t *TempSyncMsg) SetMsgChain(mc msgchain.MsgChain) {
	t.MsgChain = mc
}

// StrangerSyncMsg is synchronized StrangerMsg.
type StrangerSyncMsg struct {
	Type Type `json:"type" type:"StrangerSyncMessage"`

	RawMsgChain json.RawMessage `json:"messageChain"`
	MsgChain    msgchain.MsgChain

	// Subject is the current stranger's account to which the message was sent.
	Subject Stranger `json:"subject"`
}

func (s *StrangerSyncMsg) GetType() Type {
	return TypeStrangerSyncMsg
}

func (s *StrangerSyncMsg) GetRawMsgChain() json.RawMessage {
	return s.RawMsgChain
}

func (s *StrangerSyncMsg) SetMsgChain(mc msgchain.MsgChain) {
	s.MsgChain = mc
}

/********** BOT SELF EVENT **********/

// BotOnline indicates the bot login successful.
type BotOnline struct {
	Type Type `json:"type" type:"BotOnlineEvent"`

	// QQ number of the affected bot
	QQ int `json:"qq"`
}

func (b *BotOnline) GetType() Type {
	return TypeBotOnlineEvent
}

// BotOfflineActive indicates the bot actively offline.
type BotOfflineActive struct {
	Type Type `json:"type" type:"BotOfflineEventActive"`

	// QQ number of the affected bot
	QQ int `json:"qq"`
}

func (b *BotOfflineActive) GetType() Type {
	return TypeBotOfflineEventActive
}

// BotOfflineForce indicates the bot is forced to go offline due to logging in
// on another device with the same platform.
type BotOfflineForce struct {
	Type Type `json:"type" type:"BotOfflineEventForce"`

	// QQ number of the affected bot
	QQ int `json:"qq"`
}

func (b *BotOfflineForce) GetType() Type {
	return TypeBotOfflineEventForce
}

// BotOfflineDrop indicates the bot was disconnected from the server or dropped
// due to network problems.
type BotOfflineDrop struct {
	Type Type `json:"type" type:"BotOfflineEventDropped"`

	// QQ number of the affected bot
	QQ int `json:"qq"`
}

func (b *BotOfflineDrop) GetType() Type {
	return TypeBotOfflineEventDropped
}

// BotRelogin indicates the bot actively re-logged in.
type BotRelogin struct {
	Type Type `json:"type" type:"BotReloginEvent"`

	// QQ number of the affected bot
	QQ int `json:"qq"`
}

func (b *BotRelogin) GetType() Type {
	return TypeBotReloginEvent
}

/*********** FRIEND EVENT ***********/

// FriendInputStatusChanged indicates the friend input status has changed.
type FriendInputStatusChanged struct {
	Type Type `json:"type" type:"FriendInputStatusChangedEvent"`

	// Involved friend
	Friend Friend `json:"friend"`

	// Current inputting state, i.e. whether it is being inputted.
	Inputting bool `json:"inputting"`
}

func (f *FriendInputStatusChanged) GetType() Type {
	return TypeFriendInputStatusChangedEvent
}

// FriendNickChanged indicates a change in a friend's nickname.
type FriendNickChanged struct {
	Type Type `json:"type" type:"FriendNickChangedEvent"`

	// (The value of Friend.Nickname is not determined in this event)
	Friend Friend `json:"friend"`

	// Original nickname
	From string `json:"from"`

	// New nickname
	To string `json:"to"`
}

func (f *FriendNickChanged) GetType() Type {
	return TypeFriendNickChangedEvent
}

// FriendAdd indicates the bot has added a new user as a friend.
type FriendAdd struct {
	Type Type `json:"type" type:"FriendAddEvent"`

	// Newly added friend
	Friend Friend `json:"friend"`

	// If the stranger is added.
	// if true it corresponds to the mirai event StrangerRelationChangeEvent.Friended,
	// otherwise it is FriendAddEvent.
	Stranger bool `json:"stranger"`
}

func (f *FriendAdd) GetType() Type {
	return TypeFriendAddEvent
}

// FriendDelete indicates a friend has been deleted.
type FriendDelete struct {
	Type Type `json:"type" type:"FriendDeleteEvent"`

	// Deleted friend
	Friend Friend `json:"friend"`
}

func (f *FriendDelete) GetType() Type {
	return TypeFriendDeleteEvent
}

/*********** GROUP EVENT ************/

// BotGroupPermissionChange indicates the bot's permission within the group has been changed,
// and the operator must be the group owner.
type BotGroupPermissionChange struct {
	Type Type `json:"type" type:"BotGroupPermissionChangeEvent"`

	// Bot's original permission
	Origin GroupPermission `json:"origin"`

	// Bot's new permission
	Current GroupPermission `json:"current"`

	// Group where the bot permissions have changed
	Group Group `json:"group"`
}

func (b *BotGroupPermissionChange) GetType() Type {
	return TypeBotGroupPermissionChangeEvent
}

// BotMute indicates the bot has been muted from the group.
type BotMute struct {
	Type Type `json:"type" type:"BotMuteEvent"`

	// Duration of the mute in seconds
	DurationSeconds int `json:"durationSecond"`

	// The administrator or group owner who performed the operation
	Operator GroupMember `json:"operator"`
}

func (b *BotMute) GetType() Type {
	return TypeBotMuteEvent
}

// BotUnmute indicates the bot has been unmuted from the group.
type BotUnmute struct {
	Type Type `json:"type" type:"BotUnmuteEvent"`

	// The administrator or group owner who performed the operation
	Operator GroupMember `json:"operator"`
}

func (b *BotUnmute) GetType() Type {
	return TypeBotUnmuteEvent
}

// BotJoinGroup indicates the bot has joined a new group.
type BotJoinGroup struct {
	Type Type `json:"type" type:"BotJoinGroupEvent"`

	// Group the bot has joined
	Group Group `json:"group"`

	// Inviter when the bot is invited to the group; otherwise, it is null.
	Invitor GroupMember `json:"invitor,omitempty"`
}

func (b *BotJoinGroup) GetType() Type {
	return TypeBotJoinGroupEvent
}

// BotLeaveActive indicates the bot has actively leave the group.
type BotLeaveActive struct {
	Type Type `json:"type" type:"BotLeaveEventActive"`

	// Group the bot has leaved
	Group Group `json:"group"`
}

func (b *BotLeaveActive) GetType() Type {
	return TypeBotLeaveEventActive
}

// BotLeaveKick indicates the bot has been kicked out of a group.
type BotLeaveKick struct {
	Type Type `json:"type" type:"BotLeaveEventKick"`

	// Group from which the bot was kicked
	Group Group `json:"group"`

	// The administrator or group owner who performed the operation
	Operator GroupMember `json:"operator"`
}

func (b *BotLeaveKick) GetType() Type {
	return TypeBotLeaveEventKick
}

// BotLeaveDisband indicates the bot leave the group due to the group owner disbanding it;
// the operator is always the group owner.
type BotLeaveDisband struct {
	Type Type `json:"type" type:"BotLeaveEventDisband"`

	// Disbanded group
	Group Group `json:"group"`

	// The administrator or group owner who disbaned the group
	Operator GroupMember `json:"operator"`
}

func (b *BotLeaveDisband) GetType() Type {
	return TypeBotLeaveEventDisband
}

// GroupRecall indicates the recall of a group message.
type GroupRecall struct {
	Type Type `json:"type" type:"GroupRecallEvent"`

	// QQ number of the original message sender
	AuthorID int `json:"authorId"`

	// MessageID of the original message
	MessageID int `json:"messageId"`

	// Time of the original message sent
	Time int `json:"time"`

	// Group in which the message recall occurred
	Group Group `json:"group"`

	// The member who performed the operation, null when bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupRecall) GetType() Type {
	return TypeGroupRecallEvent
}

// FriendRecall indicates a message being recalled in a friend conversation.
type FriendRecall struct {
	Type Type `json:"type" type:"FriendRecallEvent"`

	// QQ number of the original message sender
	AuthorID int `json:"authorId"`

	// MessageID of the original message
	MessageID int `json:"messageId"`

	// Time of the original message sent
	Time int `json:"time"`

	// QQ number of the person who initiated the recall,
	// either a friend or the bot.
	Operator int `json:"operator"`
}

func (f *FriendRecall) GetType() Type {
	return TypeFriendRecallEvent
}

// Nudge indicates a "nudge" event a.k.a. "poke".
type Nudge struct {
	Type Type `json:"type" type:"NudgeEvent"`

	// QQ number of the action initiator
	FromID int `json:"fromID"`

	// Subject represents the source of the action.
	Subject struct {

		// QQ number for friends, or group number for groups of the source.
		ID int `json:"id"`

		// Type of source, either "Friend" for friends or "Group" for groups.
		Kind string `json:"kind"`
	} `json:"subject"`

	// Type of the action
	Action string `json:"action"`

	// Custom suffix of the action
	Suffix string `json:"suffix"`

	// QQ number of the action receiver
	Target int `json:"target"`
}

func (n *Nudge) GetType() Type {
	return TypeNudgeEvent
}

// GroupNameChange indicates a change in the group name.
type GroupNameChange struct {
	Type Type `json:"type" type:"GroupNameChangeEvent"`

	// Original group name
	Origin string `json:"origin"`

	// Current group name
	Current string `json:"current"`

	// Group in which the name has been changed
	Group Group `json:"group"`

	// Operator who changed the name; it is null for bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupNameChange) GetType() Type {
	return TypeGroupNameChangeEvent
}

// GroupEntAnnChange, or GroupEntranceAnnouncement,
// indicates a change in the group's entrance announcement.
type GroupEntAnnChange struct {
	Type Type `json:"type" type:"GroupEntranceAnnouncementChangeEvent"`

	// Original entrance announcement
	Origin string `json:"origin"`

	// Current entrance announcement
	Current string `json:"current"`

	// Group in which the entrance announcement has been changed
	Group Group `json:"group"`

	// Operator who changed the entrance announcement; it is null for bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupEntAnnChange) GetType() Type {
	return TypeGroupEntranceAnnouncementChangeEvent
}

// GroupMuteAll indicates a group-wide mute for all members.
type GroupMuteAll struct {
	Type Type `json:"type" type:"GroupMuteAllEvent"`

	// Original status
	Origin bool `json:"origin"`

	// Current status
	Current bool `json:"current"`

	// Affected group
	Group Group `json:"group"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupMuteAll) GetType() Type {
	return TypeGroupMuteAllEvent
}

// GroupAllowAnonChat indicates a change in the group's anonymous chat enablement status.
type GroupAllowAnonChat struct {
	Type Type `json:"type" type:"GroupAllowAnonymousChatEvent"`

	// Original status
	Origin bool `json:"origin"`

	// Current status
	Current bool `json:"current"`

	// Affected group
	Group Group `json:"group"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupAllowAnonChat) GetType() Type {
	return TypeGroupAllowAnonymousChatEvent
}

// GroupAllowConfessTalk indicates a change in the group's "confess talk" enablement status.
type GroupAllowConfessTalk struct {
	Type Type `json:"type" type:"GroupAllowConfessTalkEvent"`

	// Original status
	Origin bool `json:"origin"`

	// Current status
	Current bool `json:"current"`

	// Affected group
	Group Group `json:"group"`

	// Whether the operation was performed by bot
	IsByBot bool `json:"isByBot"`
}

func (g *GroupAllowConfessTalk) GetType() Type {
	return TypeGroupAllowConfessTalkEvent
}

// GroupAllowAnonChat indicates a change in the policy that allow members to invite friends to the group.
type GroupAllowMemberInvite struct {
	Type Type `json:"type" type:"GroupAllowMemberInviteEvent"`

	// Original status
	Origin bool `json:"origin"`

	// Current status
	Current bool `json:"current"`

	// Affected group
	Group Group `json:"group"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operator GroupMember `json:"operator,omitempty"`
}

func (g *GroupAllowMemberInvite) GetType() Type {
	return TypeGroupAllowMemberInviteEvent
}

// MemberJoin indicates a new member joining the group.
type MemberJoin struct {
	Type Type `json:"type" type:"MemberJoinEvent"`

	// Member who joined the group
	Member GroupMember `json:"member"`

	// Inviter when the member is invited to the group; otherwise, it is null.
	Invitor GroupMember `json:"invitor,omitempty"`
}

func (m *MemberJoin) GetType() Type {
	return TypeMemberJoinEvent
}

// MemberLeaveKick indicates a member being kicked out of the group.
type MemberLeaveKick struct {
	Type Type `json:"type" type:"MemberLeaveEventKick"`

	// Member kicked out of the group
	Member GroupMember `json:"member"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operator GroupMember `json:"operator"`
}

func (m *MemberLeaveKick) GetType() Type {
	return TypeMemberLeaveEventKick
}

// MemberLeaveQuit indicates a group member actively leaving the group (not the bot).
type MemberLeaveQuit struct {
	Type Type `json:"type" type:"MemberLeaveEventQuit"`

	// Member who left the group
	Member GroupMemberSimple `json:"member"`
}

func (m *MemberLeaveQuit) GetType() Type {
	return TypeMemberLeaveEventQuit
}

// MemberCardChange indicates a change in the group member's in-group nickname, or "group card".
type MemberCardChange struct {
	Type Type `json:"type" type:"MemberCardChangeEvent"`

	// Original in-group nickname
	Origin string `json:"origin"`

	// Current in-group nickname
	Current string `json:"current"`

	// Member that is affected
	Member GroupMember `json:"member"`
}

func (m *MemberCardChange) GetType() Type {
	return TypeMemberCardChangeEvent
}

// MemberSpecialTitleChange
type MemberSpecialTitleChange struct {
	Type Type `json:"type" type:"MemberSpecialTitleChangeEvent"`

	// Original special title
	Origin string `json:"origin"`

	// Current special title
	Current string `json:"current"`

	// Affected member
	Member GroupMemberSimple `json:"member"`
}

func (m *MemberSpecialTitleChange) GetType() Type {
	return TypeMemberSpecialTitleChangeEvent
}

// MemberPermissionChange
type MemberPermissionChange struct {
	Type Type `json:"type" type:"MemberPermissionChangeEvent"`

	// Original permission
	Origin GroupPermission `json:"origin"`

	// Current permission
	Current GroupPermission `json:"current"`

	// Affected member
	Member GroupMemberSimple `json:"member"`
}

func (m *MemberPermissionChange) GetType() Type {
	return TypeMemberPermissionChangeEvent
}

// MemberMute
type MemberMute struct {
	Type Type `json:"type" type:"MemberMuteEvent"`

	// Duration of the mute in seconds
	DurationSeconds int `json:"durationSeconds"`

	// Affected member
	Member GroupMember `json:"member"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operation GroupMember `json:"operation"`
}

func (m *MemberMute) GetType() Type {
	return TypeMemberMuteEvent
}

// MemberUnmute
type MemberUnmute struct {
	Type Type `json:"type" type:"MemberUnmuteEvent"`

	// Affected member
	Member GroupMember `json:"member"`

	// The administrator or group owner who performed the operation,
	// null when bot.
	Operator GroupMember `json:"operator"`
}

func (m *MemberUnmute) GetType() Type {
	return TypeMemberMuteEvent
}

// MemberHonorChange indicates a change in the group member's honor.
type MemberHonorChange struct {
	Type Type `json:"type" type:"MemberHonorChangeEvent"`

	// Affected member
	Member GroupMember `json:"member"`

	// The action related to honor change: "achieve" for gaining an honor,
	// and "lose" for losing an honor.
	Action string `json:"action"`

	// Honor name
	Honor string `json:"honor"`
}

func (m *MemberHonorChange) GetType() Type {
	return TypeMemberHonorChangeEvent
}

/********** Request Events **********/

// NewFriendRequest indicates a new friend request.
type NewFriendRequest struct {
	Type Type `json:"type" type:"NewFriendRequestEvent"`

	// Event identifier that can be used when responding to this event.
	EventID int `json:"eventId"`

	// Applicant's QQ number
	FromID int `json:"fromId"`

	// Group number if the request is made through a specific group;
	// otherwise, it is 0.
	GroupID int `json:"groupId"`

	// Applicant's nickname,
	// or in-group nickname (if the request is made through a specific group).
	Nick int `json:"nick"`

	// Applicant's verification description
	Message int `json:"message"`
}

func (n *NewFriendRequest) GetType() Type {
	return TypeNewFriendRequestEvent
}

// MemberJoinRequest indicates a new group join request.
type MemberJoinRequest struct {
	Type Type `json:"type" type:"MemberJoinRequestEvent"`

	// Event identifier that can be used when responding to this event.
	EventID int `json:"eventId"`

	// Applicant's QQ number
	FromID int `json:"fromId"`

	// Number of the group applicant is applying to join
	GroupID int `json:"groupId"`

	// Name of the group applicant is applying to join
	GroupName string `json:"groupName"`

	// Applicant's nickname
	Nick int `json:"nick"`

	// Applicant's verification description
	Message int `json:"message"`

	// QQ number of the inviter when the member is invited to the group;
	// otherwise, it is empty.
	InvitorID int `json:"invitorId,omitempty"`
}

func (m *MemberJoinRequest) GetType() Type {
	return TypeMemberJoinRequestEvent
}

// BotInvitedJoinGroupRequest indicates the bot being invited to join a group.
type BotInvitedJoinGroupRequest struct {
	Type Type `json:"type" type:"BotInvitedJoinGroupRequestEvent"`

	// Event identifier that can be used when responding to this event.
	EventID int `json:"eventId"`

	// Inviter's QQ number
	FromID int `json:"fromId"`

	// Number of the group bot was invited to join
	GroupID int `json:"groupId"`

	// Name of the group bot was invited to join
	GroupName string `json:"groupName"`

	// Inviter's nickname
	Nick int `json:"nick"`

	// Inviter's invitation description
	Message int `json:"message"`
}

func (b *BotInvitedJoinGroupRequest) GetType() Type {
	return TypeBotInvitedJoinGroupRequestEvent
}

/******** OTHER CLIENT EVENT ********/

// OtherClientOnline
// TODO
type OtherClientOnline struct {
	Type Type `json:"type" type:"OtherClientOnlineEvent"`
}

func (o *OtherClientOnline) GetType() Type {
	return TypeOtherClientOnlineEvent
}

// OtherClientOffline
// TODO
type OtherClientOffline struct {
	Type Type `json:"type" type:"OtherClientOfflineEvent"`
}

func (o *OtherClientOffline) GetType() Type {
	return TypeOtherClientOfflineEvent
}

/********** COMMAND EVENT ***********/

// CommandExecuted
// TODO
type CommandExecuted struct {
	Type Type `json:"type" type:"CommandExecutedEvent"`
}

func (c *CommandExecuted) GetType() Type {
	return TypeCommandExecutedEvent
}
