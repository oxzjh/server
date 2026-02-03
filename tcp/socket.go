package tcp

import (
	"log"
	"net"
	"time"

	"github.com/oxzjh/server"
)

type socket struct {
	server.Socket
	conn       net.Conn
	remoteAddr string
	remoteIP   string
}

func (s *socket) read(ts *tcpServer) {
	s.conn.SetReadDeadline(time.Now().Add(ts.opts.connTimeout))
	message, err := ts.opts.parser.Parse(s.conn)
	if err == nil {
		err = ts.UpgradeSocket(s, message)
	}
	if err != nil {
		log.Println("UPGRADE:", s.GetRemoteAddr(), err)
		s.Close()
		return
	}
	for {
		message, err := ts.opts.parser.Parse(s.conn)
		if err != nil {
			ts.CloseSocket(s)
			return
		}
		ts.OnMessage(s, message, ts.Send)
		s.conn.SetReadDeadline(time.Now().Add(ts.opts.timeout))
	}
}

func (s *socket) Write(data []byte) {
	if s.conn != nil {
		s.conn.Write(data)
	}
}

func (s *socket) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

func (s *socket) GetRemoteAddr() string {
	if s.remoteAddr == "" {
		s.remoteAddr = s.conn.RemoteAddr().String()
	}
	return s.remoteAddr
}

func (s *socket) GetRemoteIP() string {
	if s.remoteIP == "" {
		s.remoteIP = s.conn.RemoteAddr().(*net.TCPAddr).IP.String()
	}
	return s.remoteIP
}
