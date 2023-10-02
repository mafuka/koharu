// event 定义了客户端上报事件的数据
//
// U: 23-10-02 19:31.
package event

// EventGeneral 定义了所有上报数据中所包含的通用数据。
type EventGeneral struct {
	Time     int64  `json:"time"`      // 事件发生的 unix 时间戳
	SelfID   int64  `json:"self_id"`   // 收到事件的机器人的 QQ 号
	PostType string `json:"post_type"` // 表示该上报的类型，message | message_sent | request | notice | meta_event
}

// Sender 表示消息发送者的信息。
type Sender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 昵称
	Sex      string `json:"sex"`      // 性别，male | female | unknown
	Age      int32  `json:"age"`      // 年龄
	Card     string `json:"card"`     // （群聊）群名片或备注
	Area     string `json:"area"`     // （群聊）地区
	Level    string `json:"level"`    // （群聊）成员等级
	Role     string `json:"role"`     // （群聊）角色，owner | admin | member
	Title    string `json:"title"`    // （群聊）专属头衔
	GroupID  int64  `json:"group_id"` // （临时群消息）来源群号
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

// Anonymous 表示匿名身份信息。
type Anonymous struct {
	ID   int64  `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户昵称
	Flag string `json:"flag"` // 匿名用户 flag, 在调用禁言 API 时需要传入
}

// RequestEvent 定义了请求上报中所含的通用数据。
type RequestEvent struct {
	EventGeneral
	RequestType string `json:"request_type"` // 请求类型，friend | group
}

// NoticeEvent 定义了通知上报中所含的通用数据。
type NoticeEvent struct {
	EventGeneral
	NoticeType string `json:"notice_type"` // 通知类型
}

// MetaEvent 定义了元事件上报中所含的通用数据。
type MetaEvent struct {
	EventGeneral
	MetaEventType string `json:"meta_event_type"` // 元数据类型
}
