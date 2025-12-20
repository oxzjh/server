package webdav

import (
	"net/http"
	"time"
)

type Option func(*webDAV)

func WithAuths(auths map[string]string) Option {
	return func(o *webDAV) {
		o.auths = auths
	}
}

func WithReadonly() Option {
	return func(o *webDAV) {
		o.readonly = true
	}
}

func WithTLS(cert, key string) Option {
	return func(o *webDAV) {
		o.cert = cert
		o.key = key
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *webDAV) {
		o.timeout = timeout
	}
}

func WithPrefix(prefix string) Option {
	return func(o *webDAV) {
		o.prefix = prefix
	}
}

func WithDir(dir string) Option {
	return func(o *webDAV) {
		o.dir = dir
	}
}

func WithLogger(logger func(*http.Request, error)) Option {
	return func(o *webDAV) {
		o.logger = logger
	}
}
