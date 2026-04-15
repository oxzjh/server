package http

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/rate"
)

type httpServer struct {
	timeout      time.Duration
	domains      []string
	allowHeaders string
	maxLength    int64
	onNotFound   http.HandlerFunc
	onPanic      func(*Context, any)
	middleware   Handler
	cert         string
	key          string
	group        *rate.Group
	auth         auth.IAuth
	authIgnores  map[string]struct{}
	handlers     map[string]Handler
}

func (s *httpServer) Reg(route string, handler Handler) {
	if _, ok := s.handlers[route]; ok {
		panic("duplicate register route: " + route)
	}
	s.handlers[route] = handler
}

func (s *httpServer) Set(route string, handler Handler) {
	s.handlers[route] = handler
}

func (s *httpServer) AuthIgnore(ignores ...string) {
	if s.authIgnores == nil {
		s.authIgnores = make(map[string]struct{}, len(ignores))
	}
	for _, route := range ignores {
		s.authIgnores[route] = struct{}{}
	}
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.domains != nil {
		if len(s.domains) == 1 {
			w.Header().Set("Access-Control-Allow-Origin", s.domains[0])
		} else if origin := r.Header.Get("Origin"); origin != "" {
			for _, domain := range s.domains {
				if domain == origin {
					w.Header().Set("Access-Control-Allow-Origin", domain)
					break
				}
			}
		}
	}
	if s.allowHeaders != "" {
		w.Header().Set("Access-Control-Allow-Headers", s.allowHeaders)
	}
	if r.Method == http.MethodOptions {
		return
	}
	if r.ContentLength > s.maxLength {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}
	if handler, ok := s.handlers[r.URL.Path]; ok {
		c := &Context{Request: r}
		defer func() {
			if e := recover(); e != nil {
				s.onPanic(c, e)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
		}()
		if s.group != nil && !s.group.Allow(c.GetIP()) {
			log.Println(c.GetIP(), r.RequestURI, "LIMITED")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
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
		if s.middleware != nil {
			response = s.middleware(c)
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
		s.onNotFound(w, r)
	}
}

func (s *httpServer) Serve(l net.Listener) error {
	svr := &http.Server{Handler: s, ReadTimeout: s.timeout}
	if s.cert != "" && s.key != "" {
		fmt.Println("Serve HTTPS on", l.Addr())
		return svr.ServeTLS(l, s.cert, s.key)
	}
	fmt.Println("Serve HTTP on", l.Addr())
	return svr.Serve(l)
}

func (s *httpServer) ListenAndServe(addr string) error {
	svr := &http.Server{Addr: addr, Handler: s, ReadHeaderTimeout: s.timeout}
	if s.cert != "" && s.key != "" {
		fmt.Println("Serve HTTPS on", addr)
		return svr.ListenAndServeTLS(s.cert, s.key)
	}
	fmt.Println("Serve HTTP on", addr)
	return svr.ListenAndServe()
}

func NewServer(opts ...Option) IServer {
	s := &httpServer{
		timeout:    5 * time.Second,
		maxLength:  0xFFFF,
		onNotFound: http.NotFound,
		onPanic: func(c *Context, e any) {
			log.Println(getIP(c.Request), c.Request.RequestURI, e)
			log.Writer().Write(debug.Stack())
		},
		handlers: map[string]Handler{},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
