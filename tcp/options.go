package tcp

import (
	"time"

	"github.com/oxzjh/stream"
)

type options struct {
	connTimeout time.Duration
	timeout     time.Duration
	parser      stream.IParser
	maker       stream.IMaker
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

func WithParser(parser stream.IParser) Option {
	return func(o *options) {
		o.parser = parser
	}
}

func WithMaker(maker stream.IMaker) Option {
	return func(o *options) {
		o.maker = maker
	}
}
