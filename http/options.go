package http

import (
	"net/http"
	"time"

	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/rate"
)

type Option func(*httpServer)

func WithTimeout(timeout time.Duration) Option {
	return func(s *httpServer) {
		s.timeout = timeout
	}
}

func WithDomains(domains ...string) Option {
	return func(s *httpServer) {
		s.domains = domains
	}
}

func WithAllowHeaders(allowHeaders string) Option {
	return func(s *httpServer) {
		s.allowHeaders = allowHeaders
	}
}

func WithMaxLength(maxLength int64) Option {
	return func(s *httpServer) {
		s.maxLength = maxLength
	}
}

func WithOnNotFound(onNotFound http.HandlerFunc) Option {
	return func(s *httpServer) {
		s.onNotFound = onNotFound
	}
}

func WithOnPanic(onPanic func(*Context, any)) Option {
	return func(s *httpServer) {
		s.onPanic = onPanic
	}
}

func WithMiddleware(middleware Handler) Option {
	return func(s *httpServer) {
		s.middleware = middleware
	}
}

func WithTLS(cert, key string) Option {
	return func(s *httpServer) {
		s.cert = cert
		s.key = key
	}
}

func WithRate(limit time.Duration, burst int) Option {
	return func(s *httpServer) {
		s.group = rate.NewGroup(limit, burst)
	}
}

func WithAuth(a auth.IAuth, ignores ...string) Option {
	return func(s *httpServer) {
		s.auth = a
		s.authIgnores = make(map[string]struct{}, len(ignores))
		for _, route := range ignores {
			s.authIgnores[route] = struct{}{}
		}
	}
}
