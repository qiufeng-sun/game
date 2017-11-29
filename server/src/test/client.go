package main

import (
	"fmt"
	"io"
	"net"

	"util/logs"

	"core/net/msg"
	pb "core/net/msg/protobuf"

	. "msg/pbmsg"
)

var _ = logs.Debug

//
var g_parser = &pb.PbParser{}

//
type Client struct {
	c     net.Conn
	index int
}

func (this Client) String() string {
	return fmt.Sprintf("client %v", this.index)
}

func (this Client) Infoln(v ...interface{}) {
	msg := []interface{}{"client", this.index}
	logs.Infoln(append(msg, v...)...)
}

func (this *Client) Send(msgId EMsgId, msgData interface{}) error {
	h, b, e := g_parser.Marshal(uint32(msgId), msgData)
	if e != nil {
		return e
	}

	sz := uint32(len(h) + len(b))
	buff := msg.Uint32Bytes(sz)
	buff = append(buff, h...)
	buff = append(buff, b...)

	_, e = this.c.Write(buff)
	if e != nil {
		return e
	}

	return nil
}

func (this *Client) Recv(msgId EMsgId, m interface{}) error {
	var bsz [4]byte
	_, e := io.ReadFull(this.c, bsz[:])
	if e != nil {
		return e
	}

	sz, _ := msg.Uint32ByBytes(bsz[:])

	buff := make([]byte, int(sz))
	_, e = io.ReadFull(this.c, buff)
	if e != nil {
		logs.Panicln("recv: expect msg len=", sz, "error=", e)
	}

	id, ok := msg.ParseMsgId(buff)
	if !ok || id != int32(msgId) {
		logs.Panicln("invalid msgId! recv id:", id, "expect:", msgId, ok)
	}

	e = g_parser.Unmarshal(buff, m)
	if e != nil {
		return e
	}

	return nil
}
