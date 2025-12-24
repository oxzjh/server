package http

type Handler func(*Context) IResponse

type IServer interface {
	Reg(string, Handler)
	Set(string, Handler)
	Serve(string) error
}
