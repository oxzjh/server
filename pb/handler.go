package pb

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"

	"github.com/oxzjh/server"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	server.Handler
	OnPanic   func(server.ISocket, int16, any)
	OnSuccess func(server.ISocket, IProtobuf)
	register  map[int16]*protobuf
}

func (h *Handler) OnMessage(socket server.ISocket, message []byte) (server.IResponse, any) {
	if len(message) < 2 {
		h.OnError(socket, "no protocol")
		return nil, nil
	}

	protocol := int16(message[0])<<8 | int16(message[1])

	r, ok := h.register[protocol]
	if !ok {
		h.OnError(socket, fmt.Sprintf("unregistered protocol: %d", protocol))
		return nil, protocol
	}

	p := reflect.New(r.t).Interface().(IProtobuf)
	if err := proto.Unmarshal(message[2:], p); err != nil {
		h.OnError(socket, fmt.Sprintf("%s protocol:%d", err.Error(), protocol))
		return nil, protocol
	}

	if h.OnSuccess != nil {
		h.OnSuccess(socket, p)
	}

	defer func() {
		if err := recover(); err != nil {
			h.OnPanic(socket, protocol, err)
		}
	}()
	resp := r.f(socket, p)
	if resp != nil {
		return &response{resp}, protocol
	}
	return nil, protocol
}

func (h *Handler) Reg(p IProtobuf, f func(server.ISocket, IProtobuf) IProtobuf) {
	protocol := p.GetProtocol()
	if _, ok := h.register[protocol]; ok {
		panic(fmt.Sprintf("Duplicate register protocol: %d", protocol))
	}
	h.register[protocol] = &protobuf{reflect.TypeOf(p).Elem(), f}
}

func NewHandler() *Handler {
	return &Handler{
		OnPanic: func(socket server.ISocket, protocol int16, err any) {
			log.Println("PANIC:", socket.GetRemoteAddr(), socket.GetId(), protocol, err)
			log.Writer().Write(debug.Stack())
		},
		register: map[int16]*protobuf{},
	}
}
