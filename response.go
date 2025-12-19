package server

import "encoding/json"

type IResponse interface {
	GetData() []byte
}

type ResponseBytes []byte

func (rb ResponseBytes) GetData() []byte {
	return rb
}

type ResponseString string

func (rs ResponseString) GetData() []byte {
	return []byte(rs)
}

type responseJson struct {
	data any
}

func (rj *responseJson) GetData() []byte {
	data, _ := json.Marshal(rj.data)
	return data
}

func NewJson(data any) IResponse {
	return &responseJson{data}
}
