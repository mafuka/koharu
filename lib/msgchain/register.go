package msgchain

import "reflect"

// msgMap maps Msg types to their corresponding structures.
var msgMap = make(map[Type]reflect.Type)

// Register registers all Msg types to the msgMap so they can be parsed to.
// It actually uses reflection and structure tags.
//
// Note: Msg types should only be registered once at program startup.
//
//	 package main
//		func main() {msgchain.Register()}
func Register() {
	msgs := []Msg{
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

	for _, m := range msgs {
		t := m.GetType()
		msgMap[Type(t)] = reflect.TypeOf(m).Elem()
	}
}
