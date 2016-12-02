package event

import (
	"testing"

	"util/logs"
)

//
func TestEvent(t *testing.T) {
	type EvtLvUp struct {
		PlayerId int
		From     int
		To       int
	}

	type EvtMoneyChange struct {
		PlayerId int
		From     int
		To       int
	}
	Register(&EvtLvUp{}, func(evt interface{}) {
		logs.Infoln("EventLvUp:", 1, evt)
	})
	Register(&EvtLvUp{}, func(evt interface{}) {
		logs.Infoln("EventLvUp:", 2, evt)
	})

	Register(&EvtMoneyChange{}, func(evt interface{}) {
		logs.Infoln("EvtMoneyChange:", 1, evt)
	})
	Register(&EvtMoneyChange{}, func(evt interface{}) {
		logs.Infoln("EvtMoneyChange:", 2, evt)
	})

	Proc(&EvtLvUp{10, 1, 2})
	Proc(&EvtMoneyChange{11, 40, 100})
}
