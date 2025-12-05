package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/oxzjh/server"
)

var (
	EnableLog bool
	MaxMemory int64 = 100 << 20
)

func ParseQuery(c *Context) IResponse {
	c.Request.Form, _ = url.ParseQuery(c.Request.URL.RawQuery)
	if c.Request.Form == nil {
		c.Request.Form = make(url.Values)
	}
	if EnableLog {
		log.Println(c.Request.RequestURI)
	}
	return nil
}

func ParseJson(c *Context) IResponse {
	c.Data = server.Map{}
	if err := json.NewDecoder(c.Request.Body).Decode(&c.Data); err != nil {
		log.Println(c.GetIP(), c.Request.RequestURI, err)
		return NewStatus(http.StatusBadRequest, err.Error())
	}
	if EnableLog {
		log.Println(c.Request.RequestURI, c.Data)
	}
	return nil
}

func ParseForm(c *Context) IResponse {
	if err := c.Request.ParseForm(); err != nil {
		log.Println(c.GetIP(), c.Request.RequestURI, err)
		return NewStatus(http.StatusBadRequest, err.Error())
	}
	if EnableLog {
		log.Println(c.Request.RequestURI, c.Request.PostForm)
	}
	return nil
}

func ParseMultipart(c *Context) IResponse {
	c.Request.ParseMultipartForm(MaxMemory)
	if EnableLog {
		log.Println(c.Request.RequestURI, c.Request.PostForm)
	}
	return nil
}

func VoidParser(*Context) IResponse {
	return nil
}
