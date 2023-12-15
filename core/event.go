package core

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"reflect"
)

// Event defines the base interface for all types of events received from Mirai.
type Event interface {
	GetType() EventType
}

type MsgEvent interface {
	GetMsg() Msg
}

// eventMap maps EventType to their corresponding struct types for dynamic event type handling.
var eventMap = make(map[EventType]reflect.Type)

// RegisterEvent initializes eventMap with all supported events. This should be called once at program startup.
func RegisterEvent() {
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

// ParseEvent takes a JSON byte slice and attempts to parse it into an appropriate Event struct.
func ParseEvent(data []byte) (Event, error) {
	temp := struct {
		Type EventType `json:"type"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	t, ok := eventMap[temp.Type]
	if !ok {
		return nil, xerrors.Errorf("unknown event type: %s", temp.Type)
	}

	ptr := reflect.New(t).Interface().(Event)

	if err := json.Unmarshal(data, ptr); err != nil {
		return nil, xerrors.Errorf("jso")
	}

	return ptr, nil
}

// EventType is a string alias used for type safety in identifying event types.
type EventType string

const (
	TypeFriendMsg   EventType = "FriendMessage"
	TypeGroupMsg    EventType = "GroupMessage"
	TypeTempMsg     EventType = "TempMessage"
	TypeStrangerMsg EventType = "StrangerMessage"

	TypeOtherClientMsg  EventType = "OtherClientMessage"
	TypeFriendSyncMsg   EventType = "FriendSyncMessage"
	TypeGroupSyncMsg    EventType = "GroupSyncMessage"
	TypeTempSyncMsg     EventType = "TempSyncMessage"
	TypeStrangerSyncMsg EventType = "StrangerSyncMessage"

	TypeBotOnlineEvent         EventType = "BotOnlineEvent"
	TypeBotOfflineEventActive  EventType = "BotOfflineEventActive"
	TypeBotOfflineEventForce   EventType = "BotOfflineEventForce"
	TypeBotOfflineEventDropped EventType = "BotOfflineEventDropped"
	TypeBotReloginEvent        EventType = "BotReloginEvent"

	TypeFriendInputStatusChangedEvent EventType = "FriendInputStatusChangedEvent"
	TypeFriendNickChangedEvent        EventType = "FriendNickChangedEvent"
	TypeFriendAddEvent                EventType = "FriendAddEvent"
	TypeFriendDeleteEvent             EventType = "FriendDeleteEvent"

	TypeBotGroupPermissionChangeEvent        EventType = "BotGroupPermissionChangeEvent"
	TypeBotMuteEvent                         EventType = "BotMuteEvent"
	TypeBotUnmuteEvent                       EventType = "BotUnmuteEvent"
	TypeBotJoinGroupEvent                    EventType = "BotJoinGroupEvent"
	TypeBotLeaveEventActive                  EventType = "BotLeaveEventActive"
	TypeBotLeaveEventKick                    EventType = "BotLeaveEventKick"
	TypeBotLeaveEventDisband                 EventType = "BotLeaveEventDisband"
	TypeGroupRecallEvent                     EventType = "GroupRecallEvent"
	TypeFriendRecallEvent                    EventType = "FriendRecallEvent"
	TypeNudgeEvent                           EventType = "NudgeEvent"
	TypeGroupNameChangeEvent                 EventType = "GroupNameChangeEvent"
	TypeGroupEntranceAnnouncementChangeEvent EventType = "GroupEntranceAnnouncementChangeEvent"
	TypeGroupMuteAllEvent                    EventType = "GroupMuteAllEvent"
	TypeGroupAllowAnonymousChatEvent         EventType = "GroupAllowAnonymousChatEvent"
	TypeGroupAllowConfessTalkEvent           EventType = "GroupAllowConfessTalkEvent"
	TypeGroupAllowMemberInviteEvent          EventType = "GroupAllowMemberInviteEvent"
	TypeMemberJoinEvent                      EventType = "MemberJoinEvent"
	TypeMemberLeaveEventKick                 EventType = "MemberLeaveEventKick"
	TypeMemberLeaveEventQuit                 EventType = "MemberLeaveEventQuit"
	TypeMemberCardChangeEvent                EventType = "MemberCardChangeEvent"
	TypeMemberSpecialTitleChangeEvent        EventType = "MemberSpecialTitleChangeEvent"
	TypeMemberPermissionChangeEvent          EventType = "MemberPermissionChangeEvent"
	TypeMemberMuteEvent                      EventType = "MemberMuteEvent"
	TypeMemberUnmuteEvent                    EventType = "MemberUnmuteEvent"
	TypeMemberHonorChangeEvent               EventType = "MemberHonorChangeEvent"

	TypeNewFriendRequestEvent           EventType = "NewFriendRequestEvent"
	TypeMemberJoinRequestEvent          EventType = "MemberJoinRequestEvent"
	TypeBotInvitedJoinGroupRequestEvent EventType = "BotInvitedJoinGroupRequestEvent"

	TypeOtherClientOnlineEvent  EventType = "OtherClientOnlineEvent"
	TypeOtherClientOfflineEvent EventType = "OtherClientOfflineEvent"

	TypeCommandExecutedEvent EventType = "CommandExecutedEvent"
)

/*********** COMMON TYPE ************/
// Common types are used across different event structures for consistent data representation.

// Friend represents a friend in the context of friend-related events and messages.
type Friend struct {
	ID       int    `json:"ID"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

// GroupMember represents a member in a group, used in group-related messages.
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

// GroupMemberSimple is a simplified version of GroupMember with fewer fields.
type GroupMemberSimple struct {
	ID         int    `json:"id"`
	MemberName string `json:"memberName"`
	Permission string `json:"permission"`
	Group      Group  `json:"group"`
}

// Group represents a group chat.
type Group struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Permission GroupPermission `json:"permission"` // Bot permission
}

// GroupPermission defines the permission level of a member within a group.
type GroupPermission string

const (
	Owner  GroupPermission = "OWNER"
	Admin  GroupPermission = "ADMINISTRATOR"
	Member GroupPermission = "MEMBER"
)

// Stranger represents a stranger.
type Stranger struct {
	ID       int    `json:"ID"`
	Nickname string `json:"nickname"`

	// (? IDK why strangers would have remarks)
	Remark string `json:"remark"`
}

// OtherClient represents another client of the same account.
type OtherClient struct {
	ID       int    `json:"id"`
	Platform string `json:"platform"`
}

/********* REGULAR MESSAGE **********/
// These types represent various forms of messages received, like friend, group, or temporary messages.

// FriendMsg represents a message from a friend.
type FriendMsg struct {
	Type   EventType `json:"type"`
	Sender Friend    `json:"sender"`
	Msg    Msg       `json:"messageChain"`
}

func (f *FriendMsg) GetType() EventType {
	return TypeFriendMsg
}

func (f *FriendMsg) UnmarshalJSON(data []byte) error {
	type alias FriendMsg
	t := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(f),
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	c, err := ParseMsg(t.Msg)
	if err != nil {
		return nil
	}
	f.Msg = c

	return nil
}

// GroupMsg represents a message from a group.
type GroupMsg struct {
	Type   EventType   `json:"type"`
	Sender GroupMember `json:"sender"`
	Msg    Msg         `json:"messageChain"`
}

func (g *GroupMsg) GetType() EventType {
	return TypeGroupMsg
}

func (g *GroupMsg) UnmarshalJSON(data []byte) error {
	type alias GroupMsg
	t := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(g),
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	c, err := ParseMsg(t.Msg)
	if err != nil {
		return nil
	}
	g.Msg = c

	return nil
}

// TempMsg represents a message from a temporary chat within a group.
type TempMsg struct {
	Type   EventType   `json:"type"`
	Sender GroupMember `json:"sender"`
	Msg    Msg         `json:"messageChain"`
}

func (t *TempMsg) GetType() EventType {
	return TypeTempMsg
}

func (t *TempMsg) UnmarshalJSON(data []byte) error {
	type alias TempMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	t.Msg = c

	return nil
}

// StrangerMsg represents a message from a stranger.
type StrangerMsg struct {
	Type EventType `json:"type"`

	Sender Stranger `json:"sender"`

	Msg Msg `json:"messageChain"`
}

func (s *StrangerMsg) GetType() EventType {
	return TypeStrangerMsg
}

func (s *StrangerMsg) UnmarshalJSON(data []byte) error {
	type alias StrangerMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(s),
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	s.Msg = c

	return nil
}

// OtherClientMsg represents a message from another client of the same account.
type OtherClientMsg struct {
	Type   EventType   `json:"type"`
	Sender OtherClient `json:"sender"`
	Msg    Msg         `json:"messageChain"`
}

func (o *OtherClientMsg) GetType() EventType {
	return TypeOtherClientMsg
}

func (o *OtherClientMsg) UnmarshalJSON(data []byte) error {
	type alias OtherClientMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(o),
	}

	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	o.Msg = c

	return nil
}

/******* SYNCHRONIZED MESSAGE *******/
// Synchronized messages are events generated when messages sent by other clients are synchronized to the Mirai bot.
// The sender of such events is always the bot itself, so they are omitted.

type FriendSyncMsg struct {
	Type    EventType `json:"type"`
	Subject Friend    `json:"subject"` // the target friend for sending
	Msg     Msg       `json:"messageChain"`
}

func (f *FriendSyncMsg) GetType() EventType {
	return TypeFriendSyncMsg
}

func (f *FriendSyncMsg) UnmarshalJSON(data []byte) error {
	type alias FriendSyncMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(f),
	}

	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	f.Msg = c

	return nil
}

type GroupSyncMsg struct {
	Type    EventType `json:"type"`
	Subject Group     `json:"subject"` // the target group for sending.
	Msg     Msg       `json:"messageChain"`
}

func (g *GroupSyncMsg) GetType() EventType {
	return TypeGroupSyncMsg
}

func (g *GroupSyncMsg) UnmarshalJSON(data []byte) error {
	type alias GroupSyncMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(g),
	}

	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	g.Msg = c

	return nil
}

// TempSyncMsg is synchronized TempMsg.
type TempSyncMsg struct {
	Type    EventType   `json:"type"`
	Subject GroupMember `json:"subject"` // the target group member
	Msg     Msg         `json:"messageChain"`
}

func (t *TempSyncMsg) GetType() EventType {
	return TypeTempSyncMsg
}

func (t *TempSyncMsg) UnmarshalJSON(data []byte) error {
	type alias TempSyncMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	t.Msg = c

	return nil
}

type StrangerSyncMsg struct {
	Type    EventType `json:"type"`
	Subject Stranger  `json:"subject"` // the current stranger's account to which the message was sent
	Msg     Msg       `json:"messageChain"`
}

func (s *StrangerSyncMsg) GetType() EventType {
	return TypeStrangerSyncMsg
}

func (s *StrangerSyncMsg) UnmarshalJSON(data []byte) error {
	type alias StrangerSyncMsg
	temp := struct {
		Msg json.RawMessage `json:"messageChain"`
		*alias
	}{
		alias: (*alias)(s),
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	c, err := ParseMsg(temp.Msg)
	if err != nil {
		return nil
	}
	s.Msg = c

	return nil
}

/********** BOT SELF EVENT **********/
// These types represent various statuses and actions of the bot itself.

// BotOnline indicates the bot has logged in successfully.
type BotOnline struct {
	Type EventType `json:"type"`
	QQ   int       `json:"qq"` // QQ number of the affected bot
}

func (b *BotOnline) GetType() EventType {
	return TypeBotOnlineEvent
}

// BotOfflineActive indicates the bot has actively gone offline.
type BotOfflineActive struct {
	Type EventType `json:"type"`
	QQ   int       `json:"qq"` // QQ number of the affected bot
}

func (b *BotOfflineActive) GetType() EventType {
	return TypeBotOfflineEventActive
}

// BotOfflineForce indicates the bot was forced offline due to a simultaneous login on another device.
type BotOfflineForce struct {
	Type EventType `json:"type"`
	QQ   int       `json:"qq"` // QQ number of the affected bot
}

func (b *BotOfflineForce) GetType() EventType {
	return TypeBotOfflineEventForce
}

// BotOfflineDrop indicates the bot was disconnected due to network issues.
type BotOfflineDrop struct {
	Type EventType `json:"type"`
	QQ   int       `json:"qq"` // QQ number of the affected bot
}

func (b *BotOfflineDrop) GetType() EventType {
	return TypeBotOfflineEventDropped
}

// BotRelogin indicates the bot has actively re-logged in.
type BotRelogin struct {
	Type EventType `json:"type"`
	QQ   int       `json:"qq"` // QQ number of the affected bot
}

func (b *BotRelogin) GetType() EventType {
	return TypeBotReloginEvent
}

/*********** FRIEND EVENT ***********/
// These types represent various events related to friends, such as changes in friendship status or friend information.

// FriendInputStatusChanged indicates a change in a friend's input status.
type FriendInputStatusChanged struct {
	Type      EventType `json:"type"`
	Friend    Friend    `json:"friend"`    // Involved friend
	Inputting bool      `json:"inputting"` // Current inputting state, i.e., whether it is being inputted.
}

func (f *FriendInputStatusChanged) GetType() EventType {
	return TypeFriendInputStatusChangedEvent
}

// FriendNickChanged indicates a change in a friend's nickname.
type FriendNickChanged struct {
	Type   EventType `json:"type"`
	Friend Friend    `json:"friend"` // (The value of Friend.Nickname is not determined in this event)
	From   string    `json:"from"`   // Original nickname
	To     string    `json:"to"`     // New nickname
}

func (f *FriendNickChanged) GetType() EventType {
	return TypeFriendNickChangedEvent
}

// FriendAdd indicates the addition of a new friend.
type FriendAdd struct {
	Type   EventType `json:"type"`
	Friend Friend    `json:"friend"` // Newly added friend

	// If the stranger is added.
	// if true it corresponds to the mirai event StrangerRelationChangeEvent.Friended,
	// otherwise it is FriendAddEvent.
	// TODO: fix this comment.
	Stranger bool `json:"stranger"`
}

func (f *FriendAdd) GetType() EventType {
	return TypeFriendAddEvent
}

// FriendDelete indicates the deletion of a friend.
type FriendDelete struct {
	Type   EventType `json:"type"`
	Friend Friend    `json:"friend"` // Deleted friend
}

func (f *FriendDelete) GetType() EventType {
	return TypeFriendDeleteEvent
}

/*********** GROUP EVENT ************/
// These events represent various events related to group activities and member status within the group.

// BotGroupPermissionChange indicates a change in the bot's permission level within a group
// (the operator must be the group owner).
type BotGroupPermissionChange struct {
	Type    EventType       `json:"type"`
	Origin  GroupPermission `json:"origin"`  // Bot's original permission
	Current GroupPermission `json:"current"` // Bot's new permission
	Group   Group           `json:"group"`   // Group where the bot permissions have changed
}

func (b *BotGroupPermissionChange) GetType() EventType {
	return TypeBotGroupPermissionChangeEvent
}

// BotMute indicates the bot has been muted in a group.
type BotMute struct {
	Type            EventType   `json:"type"`
	DurationSeconds int         `json:"durationSecond"` // Duration of the mute in seconds
	Operator        GroupMember `json:"operator"`       // The administrator or group owner who performed the operation
}

func (b *BotMute) GetType() EventType {
	return TypeBotMuteEvent
}

// BotUnmute indicates the bot has been unmuted in a group.
type BotUnmute struct {
	Type     EventType   `json:"type"`
	Operator GroupMember `json:"operator"` // The administrator or group owner who performed the operation
}

func (b *BotUnmute) GetType() EventType {
	return TypeBotUnmuteEvent
}

// BotJoinGroup indicates the bot has joined a new group.
type BotJoinGroup struct {
	Type    EventType   `json:"type"`
	Group   Group       `json:"group"`             // which the bot has joined
	Invitor GroupMember `json:"invitor,omitempty"` // Inviter when the bot is invited to the group; otherwise, it is null.
}

func (b *BotJoinGroup) GetType() EventType {
	return TypeBotJoinGroupEvent
}

// BotLeaveActive indicates the bot has actively left the group.
type BotLeaveActive struct {
	Type  EventType `json:"type"`
	Group Group     `json:"group"` // which the bot has leaved
}

func (b *BotLeaveActive) GetType() EventType {
	return TypeBotLeaveEventActive
}

// BotLeaveKick indicates the bot has been kicked out of a group.
type BotLeaveKick struct {
	Type     EventType   `json:"type"`
	Group    Group       `json:"group"`    // which the bot was kicked
	Operator GroupMember `json:"operator"` // The administrator or group owner who performed the operation
}

func (b *BotLeaveKick) GetType() EventType {
	return TypeBotLeaveEventKick
}

// BotLeaveDisband indicates the bot left the group due to its disbandment
// (the operator is always the group owner).
type BotLeaveDisband struct {
	Type     EventType   `json:"type"`
	Group    Group       `json:"group"`    // Disbanded group
	Operator GroupMember `json:"operator"` // The administrator or group owner who disbanded the group
}

func (b *BotLeaveDisband) GetType() EventType {
	return TypeBotLeaveEventDisband
}

// GroupRecall indicates a message in a group has been recalled.
type GroupRecall struct {
	Type      EventType   `json:"type"`
	AuthorID  int         `json:"authorId"`           // QQ number of the original message sender
	MessageID int         `json:"messageId"`          // MessageID of the original message
	Time      int         `json:"time"`               // Time of the original message sent
	Group     Group       `json:"group"`              // Group in which the message recall occurred
	Operator  GroupMember `json:"operator,omitempty"` // The member who performed the operation, null when bot.
}

func (g *GroupRecall) GetType() EventType {
	return TypeGroupRecallEvent
}

// FriendRecall indicates a message in a friend conversation has been recalled.
type FriendRecall struct {
	Type      EventType `json:"type"`
	AuthorID  int       `json:"authorId"`  // QQ number of the original message sender
	MessageID int       `json:"messageId"` // MessageID of the original message
	Time      int       `json:"time"`      // Time of the original message sent
	Operator  int       `json:"operator"`  // QQ number of the person who initiated the recall, either a friend or the bot.
}

func (f *FriendRecall) GetType() EventType {
	return TypeFriendRecallEvent
}

// Nudge represents a "nudge" event (such as a poke).
type Nudge struct {
	Type   EventType `json:"type"`
	FromID int       `json:"fromID"` // QQ number of the action initiator

	// The Subject represents the source of the action.
	Subject struct {
		ID   int    `json:"id"`   // QQ number for friend, or group number for a group of the sources.
		Kind string `json:"kind"` // Type of source, either "Friend" for friends or "Group" for groups.
	} `json:"subject"`

	Action string `json:"action"` // Type of the action
	Suffix string `json:"suffix"` // Custom suffix of the action
	Target int    `json:"target"` // QQ number of the action receiver
}

func (n *Nudge) GetType() EventType {
	return TypeNudgeEvent
}

// GroupNameChange indicates a change in the group's name.
type GroupNameChange struct {
	Type     EventType   `json:"type"`
	Origin   string      `json:"origin"`             // Original group name
	Current  string      `json:"current"`            // Current group name
	Group    Group       `json:"group"`              // Group in which the name has been changed
	Operator GroupMember `json:"operator,omitempty"` // Operator who changed the name; it is null for bot.
}

func (g *GroupNameChange) GetType() EventType {
	return TypeGroupNameChangeEvent
}

// GroupEntAnnChange indicates a change in the group's entrance announcement.
type GroupEntAnnChange struct {
	Type     EventType   `json:"type"`
	Origin   string      `json:"origin"`             // Original entrance announcement
	Current  string      `json:"current"`            // Current entrance announcement
	Group    Group       `json:"group"`              // Group in which the entrance announcement has been changed
	Operator GroupMember `json:"operator,omitempty"` // Operator who changed the entrance announcement; it is null for bot.
}

func (g *GroupEntAnnChange) GetType() EventType {
	return TypeGroupEntranceAnnouncementChangeEvent
}

// GroupMuteAll indicates a group-wide mute for all members.
type GroupMuteAll struct {
	Type     EventType   `json:"type"`
	Origin   bool        `json:"origin"`             // Original status
	Current  bool        `json:"current"`            // Current status
	Group    Group       `json:"group"`              // Affected group
	Operator GroupMember `json:"operator,omitempty"` // The administrator or group owner who performed the operation, null when bot.
}

func (g *GroupMuteAll) GetType() EventType {
	return TypeGroupMuteAllEvent
}

// GroupAllowAnonChat indicates a change in the group's anonymous chat policy.
type GroupAllowAnonChat struct {
	Type     EventType   `json:"type"`
	Origin   bool        `json:"origin"`             // Original status
	Current  bool        `json:"current"`            // Current status
	Group    Group       `json:"group"`              // Affected group
	Operator GroupMember `json:"operator,omitempty"` // The administrator or group owner who performed the operation, null when bot.
}

func (g *GroupAllowAnonChat) GetType() EventType {
	return TypeGroupAllowAnonymousChatEvent
}

// GroupAllowConfessTalk indicates a change in the group's confess talk policy.
type GroupAllowConfessTalk struct {
	Type    EventType `json:"type"`
	Origin  bool      `json:"origin"`  // Original status
	Current bool      `json:"current"` // Current status
	Group   Group     `json:"group"`   // Affected group
	IsByBot bool      `json:"isByBot"` // Whether bot performed the operation
}

func (g *GroupAllowConfessTalk) GetType() EventType {
	return TypeGroupAllowConfessTalkEvent
}

// GroupAllowMemberInvite indicates a change in the group's member invite policy.
type GroupAllowMemberInvite struct {
	Type     EventType   `json:"type"`
	Origin   bool        `json:"origin"`             // Original status
	Current  bool        `json:"current"`            // Current status
	Group    Group       `json:"group"`              // Affected group
	Operator GroupMember `json:"operator,omitempty"` // The administrator or group owner who performed the operation, null when bot.
}

func (g *GroupAllowMemberInvite) GetType() EventType {
	return TypeGroupAllowMemberInviteEvent
}

// MemberJoin indicates a new member has joined the group.
type MemberJoin struct {
	Type    EventType   `json:"type"`
	Member  GroupMember `json:"member"`            // Member, who joined the group
	Invitor GroupMember `json:"invitor,omitempty"` // Inviter when the member is invited to the group; otherwise, it is null.
}

func (m *MemberJoin) GetType() EventType {
	return TypeMemberJoinEvent
}

// MemberLeaveKick indicates a member has been kicked out of the group.
type MemberLeaveKick struct {
	Type     EventType   `json:"type"`
	Member   GroupMember `json:"member"`   // Member kicked out of the group
	Operator GroupMember `json:"operator"` // The administrator or group owner who performed the operation, null when bot.
}

func (m *MemberLeaveKick) GetType() EventType {
	return TypeMemberLeaveEventKick
}

// MemberLeaveQuit indicates a member has voluntarily left the group.
type MemberLeaveQuit struct {
	Type   EventType         `json:"type"`
	Member GroupMemberSimple `json:"member"` // who left the group
}

func (m *MemberLeaveQuit) GetType() EventType {
	return TypeMemberLeaveEventQuit
}

// MemberCardChange indicates a change in a member's group card (nickname within the group).
type MemberCardChange struct {
	Type    EventType   `json:"type"`
	Origin  string      `json:"origin"`  // Original in-group nickname
	Current string      `json:"current"` // Current in-group nickname
	Member  GroupMember `json:"member"`  // Member that is affected
}

func (m *MemberCardChange) GetType() EventType {
	return TypeMemberCardChangeEvent
}

// MemberSpecialTitleChange indicates a change in a member's special title within the group.
type MemberSpecialTitleChange struct {
	Type    EventType         `json:"type"`
	Origin  string            `json:"origin"`  // Original special title
	Current string            `json:"current"` // Current special title
	Member  GroupMemberSimple `json:"member"`  // Affected member
}

func (m *MemberSpecialTitleChange) GetType() EventType {
	return TypeMemberSpecialTitleChangeEvent
}

// MemberPermissionChange indicates a change in a member's permission level within the group.
type MemberPermissionChange struct {
	Type    EventType         `json:"type"`
	Origin  GroupPermission   `json:"origin"`  // Original permission
	Current GroupPermission   `json:"current"` // Current permission
	Member  GroupMemberSimple `json:"member"`  // Affected member
}

func (m *MemberPermissionChange) GetType() EventType {
	return TypeMemberPermissionChangeEvent
}

// MemberMute indicates a member has been muted within the group.
type MemberMute struct {
	Type            EventType   `json:"type"`
	DurationSeconds int         `json:"durationSeconds"` // Duration of the mute in seconds
	Member          GroupMember `json:"member"`          // Affected member
	Operation       GroupMember `json:"operation"`       // The administrator or group owner who performed the operation, null when bot.
}

func (m *MemberMute) GetType() EventType {
	return TypeMemberMuteEvent
}

// MemberUnmute indicates a member has been unmuted within the group.
type MemberUnmute struct {
	Type     EventType   `json:"type"`
	Member   GroupMember `json:"member"`   // Affected member
	Operator GroupMember `json:"operator"` // The administrator or group owner who performed the operation, null when bot.
}

func (m *MemberUnmute) GetType() EventType {
	return TypeMemberUnmuteEvent
}

// MemberHonorChange indicates a change in a member's honor within the group.
type MemberHonorChange struct {
	Type   EventType   `json:"type"`
	Member GroupMember `json:"member"` // Affected member
	Action string      `json:"action"` // The action related to honor change: "achieve" for gaining an honor, and "lose" for losing an honor.
	Honor  string      `json:"honor"`  // Honor name
}

func (m *MemberHonorChange) GetType() EventType {
	return TypeMemberHonorChangeEvent
}

/********** Request Events **********/

// NewFriendRequest indicates a new friend request.
type NewFriendRequest struct {
	Type    EventType `json:"type"`
	EventID int       `json:"eventId"` // Event identifier that can be used when responding to this event.
	FromID  int       `json:"fromId"`  // Applicant's QQ number
	GroupID int       `json:"groupId"` // Group number if the request is made through a specific group; otherwise, it is 0.
	Nick    int       `json:"nick"`    // Applicant's nickname, or in-group nickname (if the request is made through a specific group).
	Message int       `json:"message"` // Applicant's verification description
}

func (n *NewFriendRequest) GetType() EventType {
	return TypeNewFriendRequestEvent
}

// MemberJoinRequest indicates a new group join request.
type MemberJoinRequest struct {
	Type      EventType `json:"type"`
	EventID   int       `json:"eventId"`             // Event identifier that can be used when responding to this event.
	FromID    int       `json:"fromId"`              // Applicant's QQ number
	GroupID   int       `json:"groupId"`             // The Number of the group applicant is applying to join
	GroupName string    `json:"groupName"`           // The Name of the group applicant is applying to join
	Nick      int       `json:"nick"`                // Applicant's nickname
	Message   int       `json:"message"`             // Applicant's verification description
	InvitorID int       `json:"invitorId,omitempty"` // QQ number of the inviter when the member is invited to the group; otherwise, it is empty.
}

func (m *MemberJoinRequest) GetType() EventType {
	return TypeMemberJoinRequestEvent
}

// BotInvitedJoinGroupRequest indicates the bot being invited to join a group.
type BotInvitedJoinGroupRequest struct {
	Type      EventType `json:"type"`
	EventID   int       `json:"eventId"`   // Event identifier that can be used when responding to this event.
	FromID    int       `json:"fromId"`    // Inviter's QQ number
	GroupID   int       `json:"groupId"`   // Number of the group bot was invited to join
	GroupName string    `json:"groupName"` // The Name of the group bot was invited to join
	Nick      int       `json:"nick"`      // Inviter's nickname
	Message   int       `json:"message"`   // Inviter's invitation description
}

func (b *BotInvitedJoinGroupRequest) GetType() EventType {
	return TypeBotInvitedJoinGroupRequestEvent
}

/******** OTHER CLIENT EVENT ********/

type OtherClientOnline struct {
	Type EventType `json:"type"`
}

func (o *OtherClientOnline) GetType() EventType {
	return TypeOtherClientOnlineEvent
}

type OtherClientOffline struct {
	Type EventType `json:"type"`
}

func (o *OtherClientOffline) GetType() EventType {
	return TypeOtherClientOfflineEvent
}

/********** COMMAND EVENT ***********/

type CommandExecuted struct {
	Type EventType `json:"type"`
}

func (c *CommandExecuted) GetType() EventType {
	return TypeCommandExecutedEvent
}
