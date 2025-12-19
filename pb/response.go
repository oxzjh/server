package pb

import (
	"google.golang.org/protobuf/proto"
)

type response struct {
	p IProtobuf
}

func (r *response) GetData() []byte {
	b, _ := proto.Marshal(r.p)
	p := r.p.GetProtocol()
	data := make([]byte, 2+len(b))
	data[0] = byte(p >> 8)
	data[1] = byte(p)
	copy(data[2:], b)
	return data
}
