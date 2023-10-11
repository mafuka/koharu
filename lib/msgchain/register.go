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
	types := []reflect.Type{
		reflect.TypeOf(Source{}),
		reflect.TypeOf(Quote{}),
		reflect.TypeOf(At{}),
		reflect.TypeOf(AtAll{}),
		reflect.TypeOf(Face{}),
		reflect.TypeOf(Plain{}),
		reflect.TypeOf(Image{}),
		reflect.TypeOf(FlashImage{}),
		reflect.TypeOf(Voice{}),
		reflect.TypeOf(XML{}),
		reflect.TypeOf(JSON{}),
		reflect.TypeOf(APP{}),
	}

	for _, t := range types {
		i := reflect.New(t).Elem()
		msgType := i.Type().Field(0).Tag.Get("type")
		msgMap[Type(msgType)] = t
	}
}
