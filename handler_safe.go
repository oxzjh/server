package server

import (
	"log"
)

type SafeHandler struct {
	Handler
	ConnectC chan ISocket
	CloseC   chan ISocket
}

func (sh *SafeHandler) OnConnect(socket ISocket) {
	log.Println("CONNECT:", socket.GetRemoteAddr(), socket.GetId())
	sh.ConnectC <- socket
}

func (sh *SafeHandler) OnClose(socket ISocket) {
	if socket.GetId() == 0 {
		log.Println("CLOSE:", socket.GetRemoteAddr(), "0")
	} else {
		sh.CloseC <- socket
	}
}
