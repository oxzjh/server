package ws

import "time"

type options struct {
	cert              string
	key               string
	readLimit         int64
	sendCap           int
	timeout           time.Duration
	readHeaderTimeout time.Duration
}

type Option func(*options)

func WithTLS(cert, key string) Option {
	return func(o *options) {
		o.cert = cert
		o.key = key
	}
}

func WithReadLimit(readLimit int64) Option {
	return func(o *options) {
		o.readLimit = readLimit
	}
}

func WithSendCap(sendCap int) Option {
	return func(o *options) {
		o.sendCap = sendCap
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.readHeaderTimeout = timeout
	}
}
