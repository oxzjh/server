package server

import (
	"github.com/oxzjh/server/auth"
)

type StreamServer struct {
	Server
	maker func(int) []byte
}

func (s *StreamServer) Send(socket ISocket, resp IResponse) {
	data := resp.GetData()
	socket.Write(s.maker(len(data)))
	socket.Write(data)
}

func (s *StreamServer) SendToAll(ids []uint64, resp IResponse) {
	data := resp.GetData()
	head := s.maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, id := range ids {
		if socket, ok := s.sockets[id]; ok {
			socket.Write(head)
			socket.Write(data)
		}
	}
}

func (s *StreamServer) SendToRoom(rid uint64, resp IResponse, except uint64) {
	data := resp.GetData()
	head := s.maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			socket.Write(head)
			socket.Write(data)
		}
	}
}

func (s *StreamServer) Broadcast(resp IResponse) {
	data := resp.GetData()
	head := s.maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.sockets {
		socket.Write(head)
		socket.Write(data)
	}
}

func NewStreamServer(a auth.IAuth, h IHandler, maker func(int) []byte) *StreamServer {
	return &StreamServer{*New(a, h), maker}
}
