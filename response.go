package server

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
