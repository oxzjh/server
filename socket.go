package server

type ISocket interface {
	Write([]byte)
	Close()
	SetId(uint64)
	GetId() uint64
	GetRid() uint64
	SetRoom(uint64, string)
	GetRoom() (uint64, string)
	GetRemoteAddr() string
	GetRemoteIP() string
}

type Socket struct {
	id    uint64
	rid   uint64
	rname string
}

func (s *Socket) SetId(id uint64) {
	s.id = id
}

func (s *Socket) GetId() uint64 {
	return s.id
}

func (s *Socket) GetRid() uint64 {
	return s.rid
}

func (s *Socket) SetRoom(id uint64, name string) {
	s.rid = id
	s.rname = name
}

func (s *Socket) GetRoom() (uint64, string) {
	return s.rid, s.rname
}
