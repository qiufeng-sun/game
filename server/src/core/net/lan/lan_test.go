package lan

import (
	"sync"
	"testing"
	"time"
)

var _ = time.Now

//
func TestLan(t *testing.T) {
	url := "tcp://127.0.0.1:8801"
	loop := 10
	wg := &sync.WaitGroup{}
	t.Log("loop:", loop)

	wg.Add(1)
	go func() {
		t.Log("create server!")
		s := NewServer(url)

		for i := 0; i < loop; i++ {
			msg, e := s.Recv()
			if e != nil {
				t.Fatal("server recv", i, e)
			}
			t.Log("server recv", i, string(msg))
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		c := NewClient(url)

		t.Log("create client!")
		for i := 0; i < loop; i++ {
			t.Log("send:", i)
			e := c.Send([]byte("test"))
			if e != nil {
				t.Fatal("client send", i, e)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
