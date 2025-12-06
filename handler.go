package server

import (
	"log"
	"time"

	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/rate"
)

type IResponse interface {
	GetData() []byte
}

type IHandler interface {
	GetServer() IServer
	SetServer(IServer)
	SetSlow(time.Duration)
	SetRate(time.Duration, int)
	Send(ISocket, IResponse)
	SendTo(uint64, IResponse)
	SendToAll([]uint64, IResponse)
	Broadcast(IResponse)
	SendToRoom(uint64, IResponse, uint64)
	ClearRoom(uint64)
	OnUpgrade([]byte, string) uint64
	OnConnect(ISocket)
	OnMessage(ISocket, []byte) IResponse
	OnLogout(ISocket)
	OnClose(ISocket)
}

type Handler struct {
	auth   auth.IAuth
	server IServer
	slow   time.Duration
	group  *rate.Group
}

func (h *Handler) GetServer() IServer {
	return h.server
}

func (h *Handler) SetServer(server IServer) {
	h.server = server
}

func (h *Handler) SetSlow(slow time.Duration) {
	h.slow = slow
}

func (h *Handler) GetSlowDuration(t time.Time) time.Duration {
	if h.slow > 0 {
		if dt := time.Since(t); dt > h.slow {
			return dt
		}
	}
	return 0
}

func (h *Handler) SetRate(limit time.Duration, burst int) {
	h.group = rate.NewGroup(limit, burst)
}

func (h *Handler) IsLimited(socket ISocket) bool {
	return h.group != nil && !h.group.Allow(socket.GetRemoteIP())
}

func (h *Handler) Send(socket ISocket, res IResponse) {
	h.server.Send(socket, res.GetData())
}

func (h *Handler) SendTo(id uint64, res IResponse) {
	if socket := h.server.GetSocket(id); socket != nil {
		h.server.Send(socket, res.GetData())
	}
}

func (h *Handler) SendToAll(ids []uint64, res IResponse) {
	h.server.SendToAll(ids, res.GetData())
}

func (h *Handler) Broadcast(res IResponse) {
	h.server.Broadcast(res.GetData())
}

func (h *Handler) SendToRoom(rid uint64, res IResponse, except uint64) {
	h.server.SendToRoom(rid, res.GetData(), except)
}

func (h *Handler) ClearRoom(rid uint64) {
	h.server.ClearRoom(rid)
}

func (h *Handler) OnUpgrade(message []byte, remoteAddr string) uint64 {
	id, err := h.auth.ParseUintToken(string(message))
	if err != nil {
		log.Println("ERROR:", remoteAddr, err)
	}
	return id
}

func (*Handler) OnConnect(socket ISocket) {
	log.Println("CONNECT:", socket.GetRemoteAddr(), socket.GetId())
}

func (*Handler) OnMessage(ISocket, []byte) IResponse {
	panic("method OnMessage not implemented")
}

func (*Handler) OnLogout(socket ISocket) {
	log.Println("LOGOUT:", socket.GetRemoteAddr(), socket.GetId())
}

func (*Handler) OnClose(socket ISocket) {
	if socket.GetId() == 0 {
		log.Println("CLOSE:", socket.GetRemoteAddr(), "0")
	}
}

func NewHandler(a auth.IAuth) *Handler {
	return &Handler{auth: a}
}
