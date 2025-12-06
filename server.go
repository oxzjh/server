package server

import "sync"

type IServer interface {
	GetSocket(uint64) ISocket
	CloseSocket(ISocket)
	GetCount() int
	GetAllUids() []uint64
	GetAllSockets() []ISocket
	EnterRoom(ISocket, uint64)
	ClearRoom(uint64)
	GetRoomSockets(uint64, uint64) []ISocket
	GetRoomUids(uint64, uint64) []uint64
	Send(ISocket, []byte)
	SendTo(uint64, []byte)
	SendToAll([]uint64, []byte)
	Broadcast([]byte)
	SendToRoom(uint64, []byte, uint64)
	Close()
}

type Server struct {
	sync.RWMutex
	Sockets map[uint64]ISocket
}

func (s *Server) GetSocket(id uint64) ISocket {
	s.RLock()
	defer s.RUnlock()
	return s.Sockets[id]
}

func (s *Server) CloseSocket(socket ISocket) {
	id := socket.GetId()
	if id > 0 {
		s.Lock()
		if s.Sockets[id] == socket {
			delete(s.Sockets, id)
		}
		s.Unlock()
	}
	socket.Close()
}

func (s *Server) GetCount() int {
	return len(s.Sockets)
}

func (s *Server) GetAllUids() []uint64 {
	uids := make([]uint64, 0, len(s.Sockets))
	s.RLock()
	for uid := range s.Sockets {
		uids = append(uids, uid)
	}
	s.RUnlock()
	return uids
}

func (s *Server) GetAllSockets() []ISocket {
	sockets := make([]ISocket, 0, len(s.Sockets))
	s.RLock()
	for _, socket := range s.Sockets {
		sockets = append(sockets, socket)
	}
	s.RUnlock()
	return sockets
}

func (s *Server) EnterRoom(socket ISocket, rid uint64) {
	socket.SetRid(rid)
}

func (s *Server) ClearRoom(rid uint64) {
	s.RLock()
	for _, socket := range s.Sockets {
		if socket.GetRid() == rid {
			socket.SetRid(0)
		}
	}
	s.RUnlock()
}

func (s *Server) GetRoomSockets(rid, except uint64) []ISocket {
	sockets := make([]ISocket, 0)
	s.RLock()
	for _, socket := range s.Sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			sockets = append(sockets, socket)
		}
	}
	s.RUnlock()
	return sockets
}

func (s *Server) GetRoomUids(rid, except uint64) []uint64 {
	uids := make([]uint64, 0)
	s.RLock()
	for _, socket := range s.Sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			uids = append(uids, socket.GetId())
		}
	}
	s.RUnlock()
	return uids
}

func (s *Server) Send(socket ISocket, data []byte) {
	socket.Write(data)
}

func (s *Server) SendTo(id uint64, data []byte) {
	if socket := s.GetSocket(id); socket != nil {
		s.Send(socket, data)
	}
}

func (s *Server) SendToAll(ids []uint64, data []byte) {
	s.RLock()
	defer s.RUnlock()
	for _, id := range ids {
		if socket, ok := s.Sockets[id]; ok {
			socket.Write(data)
		}
	}
}

func (s *Server) Broadcast(data []byte) {
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.Sockets {
		socket.Write(data)
	}
}

func (s *Server) SendToRoom(rid uint64, data []byte, except uint64) {
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.Sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			socket.Write(data)
		}
	}
}

func (s *Server) Close() {
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.Sockets {
		socket.Close()
	}
}

func (s *Server) OnMessage(handler IHandler, socket ISocket, message []byte) {
	if socket.GetId() == 0 {
		if id := handler.OnUpgrade(message, socket.GetRemoteAddr()); id > 0 {
			if old := s.GetSocket(id); old != nil {
				handler.OnLogout(old)
				old.Close()
			}
			s.Lock()
			s.Sockets[id] = socket
			s.Unlock()
			socket.SetId(id)
			handler.OnConnect(socket)
		} else {
			socket.Close()
		}
	} else {
		if res := handler.OnMessage(socket, message); res != nil {
			handler.Send(socket, res)
		}
	}
}
