package event

import (
	"reflect"

	"util"
	"util/logs"
)

var _ = logs.Debug

//
type EventFunc func(event interface{})

// event name=>[]EventFunc
var g_events = map[string][]EventFunc{}

//
func Register(event interface{}, f EventFunc) {
	name := reflect.ValueOf(event).Elem().Type().Name()
	g_events[name] = append(g_events[name], f)

	logs.Infoln("register event:", name, util.Caller(1))
}

//
func Proc(event interface{}) {
	name := reflect.ValueOf(event).Elem().Type().Name()

	logs.Debug("proc event:%v", name)

	fs, ok := g_events[name]
	if !ok {
		logs.Warn("no event function! event=%v", name)
		return
	}

	for _, f := range fs {
		f(event)
	}
}
