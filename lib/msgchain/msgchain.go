package msgchain

import (
	"encoding/json"
)

// MsgChain 消息链
//
// 它是由一个个 Msg（消息）组成的数组，是真正意义上的消息。
type MsgChain []Msg

type Msg interface {
	GetType() Type
}

// Type 消息类型
type Type string

const (
	// Metadata, always the first element of the message chain.
	TypeSource Type = "Source"

	// 引用
	TypeQuote Type = "Quote"

	// @某人
	TypeAt Type = "At"

	// @全体成员
	TypeAtAll Type = "AtAll"

	// QQ 表情
	TypeFace Type = "Face"

	// 文字
	TypePlain Type = "Plain"

	// 图片
	TypeImage Type = "Image"

	// 闪照
	TypeFlashImage Type = "FlashImage"

	// 语音
	TypeVoice Type = "Voice"

	// XML 消息
	TypeXML Type = "Xml"

	// JSON 消息
	TypeJSON Type = "Json"

	// APP 消息
	TypeAPP Type = "App"

	// 戳一戳
	TypePoke Type = "Poke"

	// 掷骰子
	TypeDice Type = "Dice"

	// 商城表情
	TypeMarketFace Type = "MarketFace"

	// 音乐分享
	TypeMusicShare Type = "MusicShare"

	// 合并转发
	TypeForwardMessage Type = "ForwardMessage"

	// 文件
	TypeFile Type = "File"

	// MiraiCode,
	// https://github.com/mamoe/mirai/blob/dev/docs/Messages.md##消息元素
	TypeMiraiCode Type = "MiraiCode"
)

// Source 元数据
//
// 其永远是消息链中的第一个元素。
type Source struct {
	Type Type `json:"type" type:"Source"`

	ID   int `json:"id"`   // 消息识别号，用于引用回复
	Time int `json:"time"` // 时间戳
}

func (s *Source) GetType() Type {
	return TypeSource
}

// Quote 引用
type Quote struct {
	Type Type `json:"type" type:"Quote"`

	// 被引用回复的原消息的 MessageID
	ID int `json:"id"`

	// 被引用回复的原消息所接收的群号，当为好友消息时为 0
	GroupID int `json:"groupID"`

	// 被引用回复的原消息的发送者的 QQ 号
	SenderID int `json:"senderID"`

	// 被引用回复的原消息的消息链对象
	Origin MsgChain `json:"origin"`
}

func (q *Quote) GetType() Type {
	return TypeQuote
}

func (q *Quote) UnmarshalJSON(data []byte) error {
	type alias Quote
	t := struct {
		Origin json.RawMessage `json:"origin"`
		*alias
	}{
		alias: (*alias)(q),
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	c, err := ParseJSON(t.Origin)
	if err != nil {
		return err
	}
	q.Origin = c

	return nil
}

// At @某人
type At struct {
	Type Type `json:"type" type:"At"`

	// 群员 QQ 号
	Target int `json:"target"`
}

func (a *At) GetType() Type {
	return TypeAt
}

// AtAll @全体成员
type AtAll struct {
	Type Type `json:"type" type:"AtAll"`
}

func (a *AtAll) GetType() Type {
	return TypeAtAll
}

// Face QQ 表情
type Face struct {
	Type Type `json:"type" type:"Face"`

	// QQ 表情编号，可选，优先高于 name
	FaceID int `json:"faceId,omitempty"`

	// QQ 表情拼音，可选
	Name string `json:"name,omitempty"`
}

func (f *Face) GetType() Type {
	return TypeFace
}

// Plain 文字
type Plain struct {
	Type Type `json:"type" type:"Plain"`

	// 文字内容
	Text string `json:"text"`
}

func (p *Plain) GetType() Type {
	return TypePlain
}

// Image 图片
type Image struct {
	Type Type `json:"type" type:"Image"`

	// 图片的 imageId，群图片与好友图片格式不同。
	// 群图片格式： "imageId": "{01E9451B-70ED-EAE3-B37C-101F1EEBF5B5}.mirai";
	// 好友图片格式： "imageId": "/f8f1ab55-bf8e-4236-b55e-955848d7069f"
	// 不为空时将忽略 url 属性。
	ImageID string `json:"imageID,omitempty"`

	// 图片的 URL，发送时可作网络图片的链接；
	// 接收时为腾讯图片服务器的链接，可用于图片下载。
	URL string `json:"url,omitempty"`

	// 图片的路径，可选。
	// 发送本地图片，路径相对于 JVM 工作路径（默认是当前路径，
	// 可通过 -Duser.dir=... 指定），也可传入绝对路径。
	Path string `json:"path,omitempty"`

	// 图片的 Base 64 编码
	Base64 string `json:"base64"`
}

func (i *Image) GetType() Type {
	return TypeImage
}

// FlashImage 闪照
type FlashImage struct {
	Type Type `json:"type" type:"FlashImage"`

	// 图片的 imageId，群图片与好友图片格式不同。
	// 不为空时将忽略 url 属性。
	ImageID string `json:"imageID,omitempty"`

	// 图片的 URL。发送时可作网络图片的链接；
	// 接收时为腾讯图片服务器的链接，可用于图片下载。
	URL string `json:"url,omitempty"`

	// 图片的路径，可选。
	// 发送本地图片，路径相对于 JVM 工作路径（默认是当前路径，
	// 可通过 -Duser.dir=... 指定），也可传入绝对路径。
	Path string `json:"path,omitempty"`

	// 图片的 Base 64 编码
	Base64 string `json:"base64"`
}

func (f *FlashImage) GetType() Type {
	return TypeFlashImage
}

type Voice struct {
	// 语音的 voiceId，不为空时将忽略 url 属性
	VoiceID string `json:"voiceId,omitempty" type:"Voice"`
	// 语音的 URL，发送时可作网络语音的链接；接收时为腾讯语音服务器的链接，可用于语音下载
	URL string `json:"url"`
}

func (v *Voice) GetType() Type {
	return TypeVoice
}

// XML 消息
type XML struct {
	Type Type   `json:"type" type:"Xml"`
	XML  string `json:"xml"` // XML 文本
}

func (x *XML) GetType() Type {
	return TypeXML
}

// JSON 消息
type JSON struct {
	Type Type   `json:"type" type:"Json"`
	JSON string `json:"json"` // JSON 文本
}

func (j *JSON) GetType() Type {
	return TypeJSON
}

// APP 消息
type APP struct {
	Type    Type   `json:"type" type:"App"`
	Content string `json:"content"` // 内容
}

func (a *APP) GetType() Type {
	return TypeAPP
}
