package ws

import (
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxzjh/server"
)

type socket struct {
	server.Socket
	conn       *websocket.Conn
	sendC      chan []byte
	remoteAddr string
	remoteIP   string
}

func (s *socket) read(svr *Server) {
	for {
		if svr.Timeout > 0 {
			s.conn.SetReadDeadline(time.Now().Add(svr.Timeout))
		}
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			svr.CloseSocket(s)
			return
		}
		svr.OnMessage(s, message)
	}
}

func (s *socket) write() {
	for data := range s.sendC {
		s.conn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func (s *socket) Write(data []byte) {
	if s.sendC == nil {
		return
	}
	if len(s.sendC) == cap(s.sendC) {
		log.Println("ERROR:", s.GetRemoteAddr(), "channel full")
		s.Close()
		return
	}
	s.sendC <- data
}

func (s *socket) Close() {
	if s.sendC != nil {
		close(s.sendC)
		s.sendC = nil
		s.conn.Close()
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
