package core

import (
	"fmt"
	"testing"
	"time"
)

//
func TestPrintPanic(t *testing.T) {
	go func() {
		defer PrintPanic()

		panic("hehe!")
	}()

	time.Sleep(time.Second)
}

//
func print0() {
	fmt.Println("no param function")
	panic("f0")
}

func print1(s string) {
	fmt.Println(s)
	panic("f1")
}

func print2(s, s2 string) {
	fmt.Println(s, s2)
	panic("f2")
}

//
func TestGoExec(t *testing.T) {
	GoExec(print0)
	GoExec(print1, "f1")
	GoExec(print2, "f2", "f22")

	time.Sleep(time.Second)
}
