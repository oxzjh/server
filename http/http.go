package http

import "net"

type Handler func(*Context) IResponse

type IServer interface {
	Reg(string, Handler)
	Set(string, Handler)
	AuthIgnore(...string)
	GenToken(uint64, uint32) string
	Serve(net.Listener) error
	ListenAndServe(string) error
}
