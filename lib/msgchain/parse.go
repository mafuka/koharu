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
	var rawChain []json.RawMessage
	if err := json.Unmarshal(data, &rawChain); err != nil {
		return nil, err
	}

	var chain MsgChain
	for _, rawMsg := range rawChain {
		temp := struct {
			Type Type `json:"type"`
		}{}

		if err := json.Unmarshal(rawMsg, &temp); err != nil {
			return nil, err
		}

		t, ok := msgMap[temp.Type]
		if !ok {
			return nil, fmt.Errorf("unknown message type: %s", temp.Type)
		}

		msg := reflect.New(t).Interface().(Msg)

		if err := json.Unmarshal(rawMsg, msg); err != nil {
			return nil, err
		}

		chain = append(chain, msg)
	}

	return chain, nil
}
