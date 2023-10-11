// event 定义了客户端上报事件的数据
package event

// EventGeneral 定义了所有上报数据中所包含的通用数据。
type EventGeneral struct {
	Time     int64  `json:"time"`      // 事件发生的 unix 时间戳
	SelfID   int64  `json:"self_id"`   // 收到事件的机器人的 QQ 号
	PostType string `json:"post_type"` // 表示该上报的类型，message | message_sent | request | notice | meta_event
}

// MessageEvent 定义了消息上报中所含的通用数据。
type MessageEvent struct {
	EventGeneral
	MessageType string    `json:"message_type"` // 消息类型，private | group
	SubType     string    `json:"sub_type"`     // 表示消息的子类型，group | public
	MessageID   int32     `json:"message_id"`   // 消息 ID
	UserID      int64     `json:"user_id"`      // 发送者 QQ 号
	Message     string    `json:"message"`      // 一个消息链
	RawMessage  string    `json:"raw_message"`  // CQ 码格式的消息
	Font        int       `json:"font"`         // 字体（可能为 0）
	Sender      Sender    `json:"sender"`       // 发送者信息
	GroupID     int64     `json:"group_id"`     // 群号
	Anonymous   Anonymous `json:"anonymous"`    // 匿名身份信息，如果不是则为 null
}

// Sender 表示消息发送者的信息。
type Sender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 昵称
	Sex      string `json:"sex"`      // 性别，male | female | unknown
	Age      int32  `json:"age"`      // 年龄
	Card     string `json:"card"`     // （仅当群聊）群名片或备注
	Area     string `json:"area"`     // （仅当群聊）地区
	Level    string `json:"level"`    // （仅当群聊）成员等级
	Role     string `json:"role"`     // （仅当群聊）角色，owner | admin | member
	Title    string `json:"title"`    // （仅当群聊）专属头衔
	GroupID  int64  `json:"group_id"` // （仅当临时群消息）来源群号
}

// Anonymous 表示匿名身份信息。
type Anonymous struct {
	ID   int64  `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户昵称
	Flag string `json:"flag"` // 匿名用户 flag, 在调用禁言 API 时需要传入
}

// RequestEvent 定义了请求上报中所含的通用数据。
type RequestEvent struct {
	EventGeneral
	RequestType string         `json:"request_type"` // 请求类型，friend | group
	UserID      int64          `json:"user_id"`      // 相关用户 QQ 号
	Comment     string         `json:"comment"`      // （仅当 friend 类型）加好友描述信息
	SubType     RequestSubType `json:"sub_type"`     // （仅当 group 类型）请求子类型，invite | ...?
	GroupID     int64          `json:"group_id"`     // （仅当 group 类型）相关群号
	InvitorID   int64          `json:"invitor_id"`   // （仅当 group 类型 invite 子类型）加群邀请人
	Flag        string         `json:"flag"`         // 请求 ID
}

type RequestSubType string

const (
	Invite RequestSubType = "invite" // 加群邀请
)

// NoticeEvent 定义了通知上报中所含的通用数据。
type NoticeEvent struct {
	EventGeneral
	NoticeType          NoticeType          `json:"notice_type"`           // 通知类型
	NoticeNotifySubType NoticeNotifySubType `json:"notice_notify_subtype"` // （通知类型为系统通知时）系统通知子类型
}

type NoticeType string

const (
	GroupUpload   NoticeType = "group_upload"   // 群文件上传
	GroupAdmin    NoticeType = "group_admin"    // 群管理员变更
	GroupDecrease NoticeType = "group_decrease" // 群成员减少
	GroupIncrease NoticeType = "group_increase" // 群成员增加
	GroupBan      NoticeType = "group_ban"      // 群成员禁言
	FriendAdd     NoticeType = "friend_add"     // 好友添加
	GroupRecall   NoticeType = "group_recall"   // 群消息撤回
	FriendRecall  NoticeType = "friend_recall"  // 好友消息撤回
	GroupCard     NoticeType = "group_card"     // 群名片变更
	OfflineFile   NoticeType = "offline_file"   // 离线文件上传
	ClientStatus  NoticeType = "client_status"  // 客户端状态变更
	Essence       NoticeType = "essence"        // 精华消息
	Notify        NoticeType = "notify"         // 系统通知
)

type NoticeNotifySubType string

const (
	Honor     NoticeNotifySubType = "honor"      // 群荣誉变更
	Poke      NoticeNotifySubType = "poke"       // 戳一戳
	LuckyKing NoticeNotifySubType = "lucky_king" // 群红包幸运王
	Title     NoticeNotifySubType = "title"      // 群成员头衔变更
)

// MetaEvent 定义了元事件上报中所含的通用数据。
type MetaEvent struct {
	EventGeneral
	MetaEventType string `json:"meta_event_type"` // 元数据类型
}
