package proxy

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

func NewReverse(targets map[string]string, errorHandler func(http.ResponseWriter, *http.Request, error)) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = targets[r.Host]
		},
		ErrorHandler: errorHandler,
	}
}

func NewRouter(defaultHost string, targets map[string]string, errorHandler func(http.ResponseWriter, *http.Request, error)) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			for route, host := range targets {
				if strings.HasPrefix(r.URL.Path, route) {
					if p := r.URL.Path[len(route):]; p == "" || p[0] == '/' {
						r.URL.Host = host
						r.URL.Path = p
						return
					}
				}
			}
			r.URL.Host = defaultHost
		},
		ErrorHandler: errorHandler,
	}
}
