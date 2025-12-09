package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/rate"
)

type Handler func(*Context) IResponse

type Server struct {
	Timeout      time.Duration
	Domains      []string
	AllowHeaders string
	MaxLength    int64
	OnNotFound   http.HandlerFunc
	OnPanic      func(*Context, any)
	Middleware   Handler
	cert         string
	key          string
	group        *rate.Group
	auth         auth.IAuth
	authIgnores  map[string]struct{}
	handlers     map[string]Handler
}

func (s *Server) SetTLS(cert, key string) {
	s.cert = cert
	s.key = key
}

func (s *Server) SetRate(limit time.Duration, burst int) {
	s.group = rate.NewGroup(limit, burst)
}

func (s *Server) SetAuth(a auth.IAuth, ignoreRoutes []string) {
	s.auth = a
	s.authIgnores = make(map[string]struct{}, len(ignoreRoutes))
	for _, route := range ignoreRoutes {
		s.authIgnores[route] = struct{}{}
	}
}

func (s *Server) Reg(route string, handler Handler) {
	if _, ok := s.handlers[route]; ok {
		panic("duplicate register route: " + route)
	}
	s.handlers[route] = handler
}

func (s *Server) Set(route string, handler Handler) {
	s.handlers[route] = handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Domains != nil {
		if len(s.Domains) == 1 {
			w.Header().Set("Access-Control-Allow-Origin", s.Domains[0])
		} else if origin := r.Header.Get("Origin"); origin != "" {
			for _, domain := range s.Domains {
				if domain == origin {
					w.Header().Set("Access-Control-Allow-Origin", domain)
					break
				}
			}
		}
	}
	if s.AllowHeaders != "" {
		w.Header().Set("Access-Control-Allow-Headers", s.AllowHeaders)
	}
	if r.Method == http.MethodOptions {
		return
	}
	if r.ContentLength > s.MaxLength {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}
	if handler, ok := s.handlers[r.URL.Path]; ok {
		c := &Context{Request: r}
		defer func() {
			if e := recover(); e != nil {
				s.OnPanic(c, e)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
		}()
		if s.group != nil && !s.group.Allow(c.GetIP()) {
			log.Println(c.GetIP(), r.RequestURI, "LIMITED")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if s.auth != nil {
			if _, ok := s.authIgnores[r.RequestURI]; !ok {
				token := r.Header.Get("token")
				if token == "" && r.Method == http.MethodGet {
					ParseQuery(c)
					c.Parser = VoidParser
					token = r.Form.Get("token")
				}
				uid, err := s.auth.ParseUintToken(token)
				if err != nil {
					log.Println(c.GetIP(), r.RequestURI, "AUTH")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				c.Uid = uid
			}
		}
		var response IResponse
		if s.Middleware != nil {
			response = s.Middleware(c)
		}
		if response == nil {
			if c.Parser == nil {
				if r.Method == http.MethodGet {
					response = ParseQuery(c)
				} else {
					contentType := r.Header.Get("Content-Type")
					if strings.HasPrefix(contentType, "application/json") {
						response = ParseJson(c)
					} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
						response = ParseForm(c)
					} else if strings.HasPrefix(contentType, "application/octet-stream") {
						response = ParseQuery(c)
					} else if strings.HasPrefix(contentType, "multipart/form-data") {
						response = ParseMultipart(c)
					}
				}
			} else {
				response = c.Parser(c)
			}
		}
		if response == nil {
			response = handler(c)
		}
		if response != nil {
			response.Write(w)
		}
	} else {
		log.Println(getIP(r), r.RequestURI, "NOT FOUND")
		s.OnNotFound(w, r)
	}
}

func (s *Server) Serve(addr string) error {
	svr := &http.Server{Addr: addr, Handler: s, ReadHeaderTimeout: s.Timeout}
	if s.cert != "" && s.key != "" {
		fmt.Println("Serve HTTPS on", addr)
		return svr.ListenAndServeTLS(s.cert, s.key)
	}
	fmt.Println("Serve HTTP on", addr)
	return svr.ListenAndServe()
}

func NewServer() *Server {
	return &Server{
		Timeout:    5 * time.Second,
		MaxLength:  0xFFFF,
		OnNotFound: http.NotFound,
		OnPanic: func(c *Context, e any) {
			log.Println(getIP(c.Request), c.Request.RequestURI, e)
			log.Writer().Write(debug.Stack())
		},
		handlers: map[string]Handler{},
	}
}
