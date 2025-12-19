package tcp

import (
	"io"
	"time"
)

type options struct {
	connTimeout time.Duration
	timeout     time.Duration
	parser      func(io.Reader) ([]byte, error)
	maker       func(int) []byte
}

type Option func(*options)

func WithConnTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.connTimeout = timeout
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func WithParser(parser func(io.Reader) ([]byte, error)) Option {
	return func(o *options) {
		o.parser = parser
	}
}

func WithMaker(maker func(int) []byte) Option {
	return func(o *options) {
		o.maker = maker
	}
}
