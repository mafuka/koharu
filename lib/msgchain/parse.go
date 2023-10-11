package msgchain

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ParseJSON accepts a JSON stream of message chain included in events,
// and returns its MsgChain structure after successful parsing.
//
// Note: There must be a Type field in the JSON,
// otherwise type matching cannot be completed.
// The target Msg type must be registered in the msgMap,
// otherwise MsgChain parsing cannot be completed.
//
//	if mc, ok := ptr.(interface{ GetMsgChain() *msgchain.MsgChain }); ok {
//		msgchain.ParseJSON()
//	}
func ParseJSON(data []byte) (MsgChain, error) {
	var r []json.RawMessage
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	var mc MsgChain
	for _, rawMessage := range r {
		var baseMsg BaseMsg
		if err := json.Unmarshal(rawMessage, &baseMsg); err != nil {
			return nil, err
		}

		msgType, ok := msgMap[baseMsg.Type]
		if !ok {
			return nil, fmt.Errorf("unknown message type: %s", baseMsg.Type)
		}

		msg := reflect.New(msgType).Interface().(Msg)
		if err := json.Unmarshal(rawMessage, msg); err != nil {
			return nil, err
		}

		mc = append(mc, msg)
	}

	return mc, nil
}
