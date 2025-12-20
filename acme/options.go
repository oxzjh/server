package acme

import (
	"net/http"
	"time"
)

type options struct {
	cacheDir     string
	renewBefore  time.Duration
	errorHandler func(http.ResponseWriter, *http.Request, error)
}

type Option func(*options)

func WithCacheDir(cacheDir string) Option {
	return func(o *options) {
		o.cacheDir = cacheDir
	}
}

func WithRenewBefore(renewBefore time.Duration) Option {
	return func(o *options) {
		o.renewBefore = renewBefore
	}
}

func WithErrorHandler(errorHandler func(http.ResponseWriter, *http.Request, error)) Option {
	return func(o *options) {
		o.errorHandler = errorHandler
	}
}
