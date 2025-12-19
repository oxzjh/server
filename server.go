package server

import (
	"log"
	"sync"
	"time"

	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/rate"
)

type IServer interface {
	GetSocket(uint64) ISocket
	CloseSocket(ISocket)
	GetCount() int
	EnterRoom(ISocket, uint64)
	ClearRoom(uint64)
	SetSlow(time.Duration)
	SetRate(time.Duration, int)
	Send(ISocket, IResponse)
	SendTo(uint64, IResponse)
	SendToAll([]uint64, IResponse)
	Broadcast(IResponse)
	SendToRoom(uint64, IResponse, uint64)
	Close()
}

type Server struct {
	sync.RWMutex
	sockets map[uint64]ISocket
	slow    time.Duration
	group   *rate.Group
	a       auth.IAuth
	h       IHandler
}

func (s *Server) GetSocket(id uint64) ISocket {
	s.RLock()
	defer s.RUnlock()
	return s.sockets[id]
}

func (s *Server) CloseSocket(socket ISocket) {
	id := socket.GetId()
	if id > 0 {
		s.Lock()
		if s.sockets[id] == socket {
			delete(s.sockets, id)
		}
		s.Unlock()
	}
	socket.Close()
	s.h.OnClose(socket)
}

func (s *Server) GetCount() int {
	return len(s.sockets)
}

func (*Server) EnterRoom(socket ISocket, rid uint64) {
	socket.SetRid(rid)
}

func (s *Server) ClearRoom(rid uint64) {
	s.RLock()
	for _, socket := range s.sockets {
		if socket.GetRid() == rid {
			socket.SetRid(0)
		}
	}
	s.RUnlock()
}

func (s *Server) SetSlow(slow time.Duration) {
	s.slow = slow
}

func (s *Server) SetRate(limit time.Duration, burst int) {
	s.group = rate.NewGroup(limit, burst)
}

func (*Server) Send(socket ISocket, resp IResponse) {
	socket.Write(resp.GetData())
}

func (s *Server) SendTo(id uint64, resp IResponse) {
	if socket := s.GetSocket(id); socket != nil {
		s.Send(socket, resp)
	}
}

func (s *Server) SendToAll(ids []uint64, resp IResponse) {
	s.RLock()
	defer s.RUnlock()
	data := resp.GetData()
	for _, id := range ids {
		if socket, ok := s.sockets[id]; ok {
			socket.Write(data)
		}
	}
}

func (s *Server) Broadcast(resp IResponse) {
	s.RLock()
	defer s.RUnlock()
	data := resp.GetData()
	for _, socket := range s.sockets {
		socket.Write(data)
	}
}

func (s *Server) SendToRoom(rid uint64, resp IResponse, except uint64) {
	s.RLock()
	defer s.RUnlock()
	data := resp.GetData()
	for _, socket := range s.sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			socket.Write(data)
		}
	}
}

func (s *Server) Close() {
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.sockets {
		socket.Close()
	}
}

func (s *Server) GetSlowDuration(t time.Time) time.Duration {
	if s.slow > 0 {
		if dt := time.Since(t); dt > s.slow {
			return dt
		}
	}
	return 0
}

func (s *Server) IsLimited(socket ISocket) bool {
	return s.group != nil && !s.group.Allow(socket.GetRemoteIP())
}

func (s *Server) Upgrade(message string) (uint64, error) {
	return s.a.ParseUintToken(message)
}

func (s *Server) AddSocket(socket ISocket, id uint64) {
	if old := s.GetSocket(id); old != nil {
		s.h.OnLogout(old)
		old.Close()
	}
	s.Lock()
	s.sockets[id] = socket
	s.Unlock()
	socket.SetId(id)
	s.h.OnConnect(socket)
}

func (s *Server) HandleMessage(socket ISocket, message []byte) {
	if resp := s.h.OnMessage(socket, message); resp != nil {
		s.Send(socket, resp)
	}
}

func (s *Server) OnMessage(socket ISocket, message []byte) {
	if socket.GetId() == 0 {
		id, err := s.Upgrade(string(message))
		if id == 0 || err != nil {
			log.Println("UPGRADE:", socket.GetRemoteAddr(), err)
			socket.Close()
			return
		}
		s.AddSocket(socket, id)
	} else {
		s.HandleMessage(socket, message)
	}
}

func New(a auth.IAuth, h IHandler) *Server {
	return &Server{
		sockets: make(map[uint64]ISocket, 1024),
		a:       a,
		h:       h,
	}
}
