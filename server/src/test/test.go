package main

import (
	"time"

	proto "github.com/golang/protobuf/proto"

	"util/logs"

	. "msg/pbmsg"
)

var _ = logs.Debug
var _ = time.Sleep

//
func TestClient(c *Client) {
	var e error

	//
	logon := &NcLogin{
		AccId:  proto.Uint32(uint32(c.index)),
		Digest: proto.String("aaa"),
		Time:   proto.Uint32(9898),
	}
	e = c.Send(EMsgId_ID_NcLogin, logon)
	checkPanic(e)

	var m NsLogin
	e = c.Recv(EMsgId_ID_NsLogin, &m)
	checkPanic(e)

	c.Infoln("nslogin:", m.String())

	//
	create := &NcCreatePlayer{Name: proto.String("nnd"), PlantIndex: proto.Int(1)}
	e = c.Send(EMsgId_ID_NcCreatePlayer, create)
	checkPanic(e)

	var screate NsCreatePlayer
	e = c.Recv(EMsgId_ID_NsCreatePlayer, &screate)
	checkPanic(e)

	c.Infoln("nscreatepalyer:", screate.String())

	//time.Sleep(time.Second * 500)
}
