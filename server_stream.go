package server

type StreamServer struct {
	Server
	Maker func(int) []byte
}

func (s *StreamServer) Send(socket ISocket, data []byte) {
	socket.Write(s.Maker(len(data)))
	socket.Write(data)
}

func (s *StreamServer) SendToAll(ids []uint64, data []byte) {
	head := s.Maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, id := range ids {
		if socket, ok := s.Sockets[id]; ok {
			socket.Write(head)
			socket.Write(data)
		}
	}
}

func (s *StreamServer) Broadcast(data []byte) {
	head := s.Maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.Sockets {
		socket.Write(head)
		socket.Write(data)
	}
}

func (s *StreamServer) SendToRoom(rid uint64, data []byte, except uint64) {
	head := s.Maker(len(data))
	s.RLock()
	defer s.RUnlock()
	for _, socket := range s.Sockets {
		if socket.GetRid() == rid && socket.GetId() != except {
			socket.Write(head)
			socket.Write(data)
		}
	}
}
