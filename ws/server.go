package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxzjh/server"
	"github.com/oxzjh/server/auth"
)

type Server struct {
	server.Server
	upgrader          *websocket.Upgrader
	cert              string
	key               string
	ReadLimit         int64
	SendCap           int
	Timeout           time.Duration
	ReadHeaderTimeout time.Duration
}

func (s *Server) SetTLS(cert, key string) {
	s.cert = cert
	s.key = key
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := s.Upgrade(r.RequestURI[1:])
	if id == 0 || err != nil {
		log.Println("UPGRADE:", r.RemoteAddr, err)
		return
	}
	if conn, err := s.upgrader.Upgrade(w, r, nil); err == nil {
		if s.ReadLimit > 0 {
			conn.SetReadLimit(s.ReadLimit)
		}
		soc := &socket{conn: conn, sendC: make(chan []byte, s.SendCap)}
		go func() {
			s.AddSocket(soc, id)
			soc.read(s)
		}()
		go soc.write()
	} else {
		log.Println("UPGRADE:", r.RemoteAddr, err)
	}
}

func (s *Server) Serve(addr string) error {
	svr := &http.Server{Addr: addr, Handler: s, ReadHeaderTimeout: s.ReadHeaderTimeout}
	if s.cert != "" && s.key != "" {
		fmt.Println("Serve WSS on", addr)
		return svr.ListenAndServeTLS(s.cert, s.key)
	}
	fmt.Println("Serve WS on", addr)
	return svr.ListenAndServe()
}

func NewServer(a auth.IAuth, h server.IHandler) *Server {
	return &Server{
		Server:            *server.New(a, h),
		upgrader:          &websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }},
		ReadLimit:         0xFFFF,
		SendCap:           128,
		Timeout:           5 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
