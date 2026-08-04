package auth

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"time"
)

type sha1Auth struct {
	secret []byte
}

func (a *sha1Auth) NewToken(data []byte, expire uint32) string {
	n := len(data)
	b := make([]byte, n+28)
	copy(b, data)
	now := uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(b[n:], now)
	binary.LittleEndian.PutUint32(b[n+4:], now+expire)
	hash := sha1.Sum(append(a.secret, b[:n+8]...))
	copy(b[n+8:], hash[:])
	return base64.URLEncoding.EncodeToString(b)
}

func (a *sha1Auth) NewUintToken(uid uint64, expire uint32) string {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, uid)
	return a.NewToken(b[:n], expire)
}

func (a *sha1Auth) ParseToken(token string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	n := len(b) - 28
	if n < 0 {
		return nil, ErrNotToken
	}
	now := uint32(time.Now().Unix())
	if binary.LittleEndian.Uint32(b[n:]) > now || binary.LittleEndian.Uint32(b[n+4:]) < now {
		return nil, ErrExpired
	}
	hash := sha1.Sum(append(a.secret, b[:n+8]...))
	if !bytes.Equal(b[n+8:], hash[:]) {
		return nil, ErrSignError
	}
	return b[:n], nil
}

func (a *sha1Auth) ParseUintToken(token string) (uint64, error) {
	b, err := a.ParseToken(token)
	if err != nil {
		return 0, err
	}
	uid, _ := binary.Uvarint(b)
	return uid, nil
}

func NewSha1(secret []byte) IAuth {
	return &sha1Auth{secret}
}
