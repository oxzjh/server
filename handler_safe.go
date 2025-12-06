package server

import (
	"log"

	"github.com/oxzjh/server/auth"
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

func NewSafeHandler(a auth.IAuth, connectC, closeC chan ISocket) *SafeHandler {
	return &SafeHandler{Handler{auth: a}, connectC, closeC}
}
