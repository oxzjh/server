package http

import (
	"net/http"
	"strings"

	"github.com/oxzjh/server"
)

type Context struct {
	Uid     uint64
	User    string
	Data    server.Map
	Request *http.Request
	Parser  Handler
	ip      string
}

func (c *Context) GetIP() string {
	if c.ip == "" {
		c.ip = getIP(c.Request)
	}
	return c.ip
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	for _, ip := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip != "" {
			return ip
		}
	}
	return r.RemoteAddr[:strings.IndexByte(r.RemoteAddr, ':')]
}
