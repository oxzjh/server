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

func NewToken(data []byte, expire uint32, secret []byte) string {
	n := len(data)
	b := make([]byte, n+12)
	copy(b, data)
	now := uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(b[n:], now)
	binary.LittleEndian.PutUint32(b[n+4:], now+expire)
	binary.LittleEndian.PutUint32(b[n+8:], crc32.ChecksumIEEE(append(secret, b[:n+8]...)))
	return base64.URLEncoding.EncodeToString(b)
}

func ParseToken(token string, secret []byte) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	n := len(b) - 12
	if n < 0 {
		return nil, ErrNotToken
	}
	now := uint32(time.Now().Unix())
	if binary.LittleEndian.Uint32(b[n:]) > now || binary.LittleEndian.Uint32(b[n+4:]) < now {
		return nil, ErrExpired
	}
	if crc32.ChecksumIEEE(append(secret, b[:n+8]...)) != binary.LittleEndian.Uint32(b[n+8:]) {
		return nil, ErrSignError
	}
	return b[:n], nil
}

func NewUintToken(uid uint64, expire uint32, secret []byte) string {
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
