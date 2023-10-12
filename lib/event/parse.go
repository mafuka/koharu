// event is a crucial component of Naki.
//
// Mirai uses Webhooks to proactively report events related to itself to the bot,
// such as changes in bot status, new friend messages, friend request additions, etc.
// These are collectively referred to as Events.
//
// For specific types of events see event.Type.
package event

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ParseJSON accepts a JSON stream of events come from Mirai's hook reporting,
// and returns its event structure after successful parsing.
//
// Note: There must be a Type field in the JSON,
// otherwise type matching cannot be completed.
// The target event must be registered in the eventMap,
// otherwise event parsing cannot be completed.
//
//	func(c *gin.Context) {
//		r, _ := c.GetRawData()
//		event.ParseJSON(r)
//	}
func ParseJSON(data []byte) (interface{}, error) {
	temp := struct {
		Type Type `json:"type"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	t, ok := eventMap[temp.Type]
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", temp.Type)
	}

	ptr := reflect.New(t).Interface()

	if err := json.Unmarshal(data, ptr); err != nil {
		return nil, err
	}

	return ptr, nil
}
