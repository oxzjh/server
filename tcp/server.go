package tcp

import (
	"net"
	"time"

	"github.com/oxzjh/server"
	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/stream"
)

type tcpServer struct {
	server.StreamServer
	opts     *options
	listener net.Listener
}

func (ts *tcpServer) Serve(addr string) (err error) {
	if ts.listener, err = net.Listen("tcp", addr); err != nil {
		return
	}
	for {
		var conn net.Conn
		if conn, err = ts.listener.Accept(); err != nil {
			return
		}
		go (&socket{conn: conn}).read(ts)
	}
}

func (ts *tcpServer) Close() {
	if ts.listener != nil {
		ts.listener.Close()
	}
	ts.StreamServer.Close()
}

func NewServer(a auth.IAuth, h server.IHandler, opts ...Option) server.IServer {
	os := &options{
		connTimeout: 5 * time.Second,
		timeout:     5 * time.Minute,
		parser:      stream.NewParser(0x92),
		maker:       stream.NewSimpleMaker(4),
	}
	for _, opt := range opts {
		opt(os)
	}
	return &tcpServer{
		StreamServer: *server.NewStreamServer(a, h, os.maker),
		opts:         os,
	}
}
