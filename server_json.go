package server

type IJsonRouter interface {
	IHandler
	Reg(string, func(ISocket, Map) IJsonResponse)
}

type jsonRequest struct {
	Method string `json:"method"`
	Data   Map    `json:"data,omitempty"`
}
