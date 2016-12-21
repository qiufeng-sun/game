package lan

import (
	"util/logs"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

var _ = logs.Debug

//
type Server struct {
	mangos.Socket
}

func NewServer(addr string) *Server {
	//
	sock, _ := pull.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if e := sock.Listen(addr); e != nil {
		logs.Panicln(e)
	}

	return &Server{Socket: sock}
}

func (this *Server) Recv() ([]byte, error) {
	return this.Socket.Recv()
}

func (this *Server) Close() {
	this.Socket.Close()
}

//
type Client struct {
	mangos.Socket
}

func NewClient(addr string) *Client {
	//
	sock, _ := push.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if e := sock.Dial(addr); e != nil {
		logs.Panicln(e)
	}

	return &Client{Socket: sock}
}

func (this *Client) Send(msg []byte) error {
	return this.Socket.Send(msg)
}

func (this *Client) Close() {
	this.Socket.Close()
}
