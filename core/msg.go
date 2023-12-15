package core

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Msg []MsgElem

type MsgElem interface {
	GetType() MsgElemType
}

// elemMap maps MsgElemType to their corresponding struct types for dynamic message parse.
var elemMap = make(map[MsgElemType]reflect.Type)

// RegisterMsgElem initializes elemMap with all supported message elements. This should be called once at program startup.
func RegisterMsgElem() {
	elem := []MsgElem{
		&Source{},
		&Quote{},
		&At{},
		&AtAll{},
		&Face{},
		&Plain{},
		&Image{},
		&FlashImage{},
		&Voice{},
		&XML{},
		&JSON{},
		&APP{},
	}

	for _, e := range elem {
		t := e.GetType()
		elemMap[t] = reflect.TypeOf(e).Elem()
	}
}

// ParseMsg takes a JSON byte slice and attempts to parse it into an appropriate MsgElem struct.
func ParseMsg(data []byte) (Msg, error) {
	var rawMsg []json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return nil, err
	}

	var msg Msg
	for _, rawMsgElem := range rawMsg {
		temp := struct {
			Type MsgElemType `json:"type"`
		}{}

		if err := json.Unmarshal(rawMsgElem, &temp); err != nil {
			return nil, err
		}

		t, ok := elemMap[temp.Type]
		if !ok {
			return nil, fmt.Errorf("unknown message type: %s", temp.Type)
		}

		element := reflect.New(t).Interface().(MsgElem)

		if err := json.Unmarshal(rawMsgElem, element); err != nil {
			return nil, err
		}

		msg = append(msg, element)
	}

	return msg, nil
}

type MsgElemType string

const (
	TypeSource         MsgElemType = "Source" // always the first element of a message.
	TypeQuote          MsgElemType = "Quote"  // reply
	TypeAt             MsgElemType = "At"     // @someone
	TypeAtAll          MsgElemType = "AtAll"  // @all
	TypeFace           MsgElemType = "Face"   // QQ emoji
	TypePlain          MsgElemType = "Plain"  // plaintext
	TypeImage          MsgElemType = "Image"
	TypeFlashImage     MsgElemType = "FlashImage"
	TypeVoice          MsgElemType = "Voice"
	TypeXML            MsgElemType = "Xml"
	TypeJSON           MsgElemType = "Json"
	TypeAPP            MsgElemType = "App"
	TypePoke           MsgElemType = "Poke"
	TypeDice           MsgElemType = "Dice"
	TypeMarketFace     MsgElemType = "MarketFace"
	TypeMusicShare     MsgElemType = "MusicShare"
	TypeForwardMessage MsgElemType = "ForwardMessage" // merged forwarding
	TypeFile           MsgElemType = "File"
	TypeMiraiCode      MsgElemType = "MiraiCode" // see https://github.com/mamoe/mirai/blob/dev/docs/Messages.md##消息元素
)

// Source represents the metadata of a message,
// and is always the first element in the message.
type Source struct {
	Type EventType `json:"type"`
	ID   int       `json:"id"`   // Message ID for referencing.
	Time int       `json:"time"` // Timestamp
}

func (s *Source) GetType() MsgElemType {
	return TypeSource
}

func (s *Source) String() string {
	return fmt.Sprintf("{ID:%d}", s.ID)
}

type Quote struct {
	Type EventType `json:"type"`

	// 被引用回复的原消息的 MessageID
	ID int `json:"id"`

	// 被引用回复的原消息所接收的群号，当为好友消息时为 0
	GroupID int `json:"groupID"`

	// 被引用回复的原消息的发送者的 QQ 号
	SenderID int `json:"senderID"`

	// 被引用回复的原消息的消息链对象
	Origin Msg `json:"origin"`
}

func (q *Quote) GetType() MsgElemType {
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

	c, err := ParseMsg(t.Origin)
	if err != nil {
		return err
	}
	q.Origin = c

	return nil
}

// At @someone
type At struct {
	Type EventType `json:"type"`

	// 群员 QQ 号
	Target int `json:"target"`
}

func (a *At) GetType() MsgElemType {
	return TypeAt
}

// AtAll @all
type AtAll struct {
	Type EventType `json:"type"`
}

func (a *AtAll) GetType() MsgElemType {
	return TypeAtAll
}

// Face QQ Emoji
type Face struct {
	Type EventType `json:"type"`

	// QQ 表情编号，可选，优先高于 name
	FaceID int `json:"faceId,omitempty"`

	// QQ 表情拼音，可选
	Name string `json:"name,omitempty"`
}

func (f *Face) GetType() MsgElemType {
	return TypeFace
}

// Plain Plaintext
type Plain struct {
	Type EventType `json:"type"`
	Text string    `json:"text"` // Content
}

func (p *Plain) GetType() MsgElemType {
	return TypePlain
}

func (p *Plain) String() string {
	return fmt.Sprintf("{Text:%s}", p.Text)
}

type Image struct {
	Type EventType `json:"type"`

	// ID of the picture, with different formats for group pictures and friend pictures.
	// Format for group pictures: "imageId": "{01E9451B-70ED-EAE3-B37C-101F1EEBF5B5}.mirai";
	// Format for friend pictures: "imageId": "/f8f1ab55-bf8e-4236-b55e-955848d7069f"
	// When not empty, the url attribute will be ignored.
	ImageID string `json:"imageID,omitempty"`

	// URL of the picture. When sending, can be used as the link for network pictures;
	// When receiving, it's the link to Tencent's image server, usable for image download.
	URL string `json:"url,omitempty"`

	// Path of the picture, optional.
	// For sending local pictures, the path is relative to the JVM working path
	// (default is the current path, can be specified by -Duser.dir=...),
	// absolute paths can also be used.
	Path string `json:"path,omitempty"`

	// Base 64 encoding of the picture
	Base64 string `json:"base64"`
}

func (i *Image) GetType() MsgElemType {
	return TypeImage
}

func (i *Image) String() string {
	return fmt.Sprintf("{Image:%s}", i.URL)
}

type FlashImage struct {
	Type EventType `json:"type"`

	// ID of the picture, different formats for group pictures and friend pictures.
	// When not empty, the url attribute will be ignored.
	ImageID string `json:"imageID,omitempty"`

	// URL of the picture. When sending, can be used as the link for network pictures;
	// When receiving, it's the link to Tencent's image server, usable for image download.
	URL string `json:"url,omitempty"`

	// Path of the picture, optional.
	// For sending local pictures, the path is relative to the JVM working path
	// (default is the current path, can be specified by -Duser.dir=...);
	// absolute paths can also be used.
	Path string `json:"path,omitempty"`

	// Base 64 encoding of the picture
	Base64 string `json:"base64"`
}

func (f *FlashImage) GetType() MsgElemType {
	return TypeFlashImage
}

type Voice struct {
	// ID of the voice; when not empty, the url attribute will be ignored.
	VoiceID string `json:"voiceId,omitempty"`

	// URL of the voice; can be used as the link for network voice when sending;
	// when receiving, it's the link to Tencent voice server, usable for voice download.
	URL string `json:"url"`
}

func (v *Voice) GetType() MsgElemType {
	return TypeVoice
}

type XML struct {
	Type EventType `json:"type"`
	XML  string    `json:"xml"`
}

func (x *XML) GetType() MsgElemType {
	return TypeXML
}

type JSON struct {
	Type EventType `json:"type"`
	JSON string    `json:"json"`
}

func (j *JSON) GetType() MsgElemType {
	return TypeJSON
}

type APP struct {
	Type    EventType `json:"type"`
	Content string    `json:"content"`
}

func (a *APP) GetType() MsgElemType {
	return TypeAPP
}
