package pb

import (
	"reflect"

	"github.com/oxzjh/server"
	"google.golang.org/protobuf/proto"
)

type IProtobufRouter interface {
	server.IServer
	Reg(IProtobuf, func(server.ISocket, IProtobuf) IProtobuf)
}

type IProtobuf interface {
	proto.Message
	GetProtocol() int16
}

type protobuf struct {
	t reflect.Type
	f func(server.ISocket, IProtobuf) IProtobuf
}
