package server

import (
	"log"
)

type IHandler interface {
	OnConnect(ISocket)
	OnMessage(ISocket, []byte) (IResponse, any)
	OnLogout(ISocket)
	OnClose(ISocket)
	OnError(ISocket, string)
}

type Handler struct{}

func (*Handler) OnConnect(socket ISocket) {
	log.Println("CONNECT:", socket.GetRemoteAddr(), socket.GetId())
}

func (*Handler) OnLogout(socket ISocket) {
	log.Println("LOGOUT:", socket.GetRemoteAddr(), socket.GetId())
}

func (*Handler) OnClose(socket ISocket) {
	if socket.GetId() == 0 {
		log.Println("CLOSE:", socket.GetRemoteAddr(), "0")
	}
}

func (*Handler) OnError(socket ISocket, err string) {
	log.Println("ERROR:", socket.GetRemoteAddr(), socket.GetId(), err)
}
