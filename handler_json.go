package server

import (
	"encoding/json"
	"log"
	"runtime/debug"
)

type JsonHandler struct {
	Handler
	OnPanic    func(ISocket, string, any)
	OnRequest  func(ISocket, string, Map)
	OnResponse func(ISocket, string, IJsonResponse)
	register   map[string]func(ISocket, Map) IJsonResponse
}

func (jh *JsonHandler) OnMessage(socket ISocket, message []byte) (IResponse, any) {
	var req jsonRequest
	if json.Unmarshal(message, &req) != nil {
		jh.OnError(socket, "unmarshal failed")
		return nil, nil
	}

	if req.Method == "" {
		jh.OnError(socket, "no method")
		return nil, nil
	}

	h, ok := jh.register[req.Method]
	if !ok {
		jh.OnError(socket, "unregistered method: "+req.Method)
		return nil, req.Method
	}

	if jh.OnRequest != nil {
		jh.OnRequest(socket, req.Method, req.Data)
	}

	defer func() {
		if err := recover(); err != nil {
			jh.OnPanic(socket, req.Method, err)
		}
	}()
	resp := h(socket, req.Data)
	if jh.OnResponse != nil {
		jh.OnResponse(socket, req.Method, resp)
	}
	if resp != nil {
		resp.SetMethod(req.Method)
	}
	return resp, req.Method
}

func (jh *JsonHandler) Reg(method string, h func(ISocket, Map) IJsonResponse) {
	if _, ok := jh.register[method]; ok {
		panic("Duplicate register method: " + method)
	}
	jh.register[method] = h
}

func NewJsonHandler() *JsonHandler {
	return &JsonHandler{
		OnPanic: func(socket ISocket, method string, err any) {
			log.Println("PANIC:", socket.GetRemoteIP(), socket.GetId(), err)
			log.Writer().Write(debug.Stack())
		},
		register: map[string]func(ISocket, Map) IJsonResponse{},
	}
}
