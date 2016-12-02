package server

import (
	"os"
	"os/signal"
	"syscall"

	"core"
)

//
func watchSignal(rch chan<- string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM:
			rch <- "signal:terminated"

		case syscall.SIGINT:
			rch <- "signal:interrupt"
		}

		close(rch)
		return
	}
}

//
func WatchSignal() <-chan string {
	c := make(chan string, 1)

	core.GoExec(watchSignal, c)

	return c
}
