package server

import "encoding/json"

type IJsonResponse interface {
	IResponse
	SetMethod(string)
}

type jsonResponse struct {
	Method string `json:"method"`
	Code   int    `json:"code,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (jr *jsonResponse) GetData() []byte {
	data, _ := json.Marshal(jr)
	return data
}

func (jr *jsonResponse) SetMethod(method string) {
	jr.Method = method
}

func NewJsonCode(code int) IJsonResponse {
	return &jsonResponse{Code: code}
}

func NewJsonData(data any) IJsonResponse {
	return &jsonResponse{Data: data}
}

func NewJsonDefault() IJsonResponse {
	return &jsonResponse{}
}

func NewJsonResponse(method string, data any) IJsonResponse {
	return &jsonResponse{Method: method, Data: data}
}
