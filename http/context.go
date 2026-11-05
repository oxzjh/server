package http

import (
	"net/http"

	"github.com/oxzjh/server"
	"github.com/oxzjh/server/ip"
)

type Context struct {
	Uid      uint64
	User     string
	Data     server.Map
	Request  *http.Request
	Parser   Handler
	Extra    any
	ipParser *ip.Parser
	ip       string
}

func (c *Context) GetIP() string {
	if c.ip == "" {
		c.ip = c.ipParser.Parse(c.Request)
	}
	return c.ip
}
