package auth

type IAuth interface {
	NewToken(data []byte, expire uint32) string
	ParseToken(token string) ([]byte, error)
	NewUintToken(uid uint64, expire uint32) string
	ParseUintToken(token string) (uint64, error)
}

type auth struct {
	secret []byte
}

func (a *auth) NewToken(data []byte, expire uint32) string {
	return NewToken(data, expire, a.secret)
}

func (a *auth) ParseToken(token string) ([]byte, error) {
	return ParseToken(token, a.secret)
}

func (a *auth) NewUintToken(uid uint64, expire uint32) string {
	return NewUintToken(uid, expire, a.secret)
}

func (a *auth) ParseUintToken(token string) (uint64, error) {
	return ParseUintToken(token, a.secret)
}

func New(secret []byte) IAuth {
	return &auth{secret}
}
