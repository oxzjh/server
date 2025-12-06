package server

import (
	"errors"
	"io"
)

var StreamKey uint16

func ParseStream(r io.Reader) ([]byte, error) {
	h := make([]byte, 4)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	n := uint16(h[0])<<8 | uint16(h[1])
	s := uint16(h[2])<<8 | uint16(h[3])
	if n^StreamKey != s {
		return nil, errors.New("sign error")
	}
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}

func ParseStream2(r io.Reader) ([]byte, error) {
	h := make([]byte, 2)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	n := uint16(h[0])<<8 | uint16(h[1])
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}

func ParseStream3(r io.Reader) ([]byte, error) {
	h := make([]byte, 3)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	n := uint32(h[0])<<16 | uint32(h[1])<<8 | uint32(h[2])
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}

func ParseStream4(r io.Reader) ([]byte, error) {
	h := make([]byte, 4)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	n := uint32(h[0])<<24 | uint32(h[1])<<16 | uint32(h[2])<<8 | uint32(h[3])
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}

func MakeStream(n int) []byte {
	s := uint16(n) ^ StreamKey
	b := make([]byte, 4)
	b[0] = byte(n >> 8)
	b[1] = byte(n)
	b[2] = byte(s >> 8)
	b[3] = byte(s)
	return b
}

func MakeStream2(n int) []byte {
	b := make([]byte, 2)
	b[0] = byte(n >> 8)
	b[1] = byte(n)
	return b
}

func MakeStream3(n int) []byte {
	b := make([]byte, 3)
	b[0] = byte(n >> 16)
	b[1] = byte(n >> 8)
	b[2] = byte(n)
	return b
}

func MakeStream4(n int) []byte {
	b := make([]byte, 4)
	b[0] = byte(n >> 24)
	b[1] = byte(n >> 16)
	b[2] = byte(n >> 8)
	b[3] = byte(n)
	return b
}
