package auth

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"time"
)

var (
	ErrNotToken  = errors.New("not token")
	ErrExpired   = errors.New("token expired")
	ErrSignError = errors.New("sign error")
)

func NewToken(data []byte, expire time.Duration, secret []byte) string {
	n := len(data)
	b := make([]byte, n+8)
	copy(b, data)
	binary.BigEndian.PutUint32(b[n:], uint32(time.Now().Add(expire).Unix()))
	binary.BigEndian.PutUint32(b[n+4:], crc32.ChecksumIEEE(append(secret, b[:n+4]...)))
	return base64.URLEncoding.EncodeToString(b)
}

func ParseToken(token string, secret []byte) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	n := len(b) - 8
	if n < 0 {
		return nil, ErrNotToken
	}
	if binary.BigEndian.Uint32(b[n:]) < uint32(time.Now().Unix()) {
		return nil, ErrExpired
	}
	if crc32.ChecksumIEEE(append(secret, b[:n+4]...)) != binary.BigEndian.Uint32(b[n+4:]) {
		return nil, ErrSignError
	}
	return b[:n], nil
}

func NewUintToken(uid uint64, expire time.Duration, secret []byte) string {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, uid)
	return NewToken(b[:n], expire, secret)
}

func ParseUintToken(token string, secret []byte) (uint64, error) {
	b, err := ParseToken(token, secret)
	if err != nil {
		return 0, err
	}
	uid, _ := binary.Uvarint(b)
	return uid, nil
}
