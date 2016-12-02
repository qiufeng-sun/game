package server

import (
	"testing"

	"time"
)

//
func TestWatchSignal(t *testing.T) {
	ch := WatchSignal()
	str := <-ch

	t.Log("get", str)
	time.Sleep(time.Second)
	t.Log("handle signal ok!!")
}
