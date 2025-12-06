package server

type ISocket interface {
	Write([]byte)
	Close()
	SetId(uint64)
	GetId() uint64
	SetRid(uint64)
	GetRid() uint64
	GetRemoteAddr() string
	GetRemoteIP() string
}

type Socket struct {
	id  uint64
	rid uint64
}

func (s *Socket) SetId(id uint64) {
	s.id = id
}

func (s *Socket) GetId() uint64 {
	return s.id
}

func (s *Socket) SetRid(rid uint64) {
	s.rid = rid
}

func (s *Socket) GetRid() uint64 {
	return s.rid
}
