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
	sh.ConnectC <- socket
}

func (sh *SafeHandler) OnClose(socket ISocket) {
	if socket.GetId() == 0 {
		log.Println("CLOSE:", socket.GetRemoteIP(), "0")
	} else {
		sh.CloseC <- socket
	}
}
