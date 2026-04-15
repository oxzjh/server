package http

import "net"

type Handler func(*Context) IResponse

type IServer interface {
	Reg(string, Handler)
	Set(string, Handler)
	AuthIgnore(...string)
	Serve(net.Listener) error
	ListenAndServe(string) error
}
