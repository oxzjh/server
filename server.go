package server

import (
	"errors"
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
	ClearRoom(uint64)
	SetSlow(time.Duration)
	SetRate(time.Duration, int)
	Send(ISocket, IResponse)
	SendTo(uint64, IResponse)
	SendToAll([]uint64, IResponse)
	Broadcast(IResponse)
	SendToRoom(uint64, IResponse, uint64)
	Close()
	Serve(string) error
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

func (s *Server) ClearRoom(rid uint64) {
	s.RLock()
	for _, socket := range s.sockets {
		if socket.GetRid() == rid {
			socket.SetRoom(0, "")
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

func (s *Server) Upgrade(message string) (uint64, error) {
	return s.a.ParseUintToken(message)
}

func (s *Server) UpgradeSocket(socket ISocket, message []byte) error {
	id, err := s.Upgrade(string(message))
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("invalid id")
	}
	s.AddSocket(socket, id)
	return nil
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

func (s *Server) OnMessage(socket ISocket, message []byte, sender func(ISocket, IResponse)) {
	if s.group != nil && !s.group.Allow(socket.GetRemoteIP()) {
		s.h.OnError(socket, "rate limited")
		return
	}

	t := time.Now()
	resp, method := s.h.OnMessage(socket, message)
	if s.slow > 0 {
		if dt := time.Since(t); dt > s.slow {
			log.Println("SLOW:", method, dt)
		}
	}
	if resp != nil {
		sender(socket, resp)
	}
}

func New(a auth.IAuth, h IHandler) *Server {
	return &Server{
		sockets: make(map[uint64]ISocket, 1024),
		a:       a,
		h:       h,
	}
}
