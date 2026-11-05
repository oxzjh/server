package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxzjh/server"
	"github.com/oxzjh/server/auth"
	"github.com/oxzjh/server/ip"
)

type wsServer struct {
	server.Server
	upgrader *websocket.Upgrader
	opts     *options
}

func (ws *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteIP := ws.opts.ipParser.Parse(r)
	id, err := ws.Upgrade(r.RequestURI[1:])
	if id == 0 || err != nil {
		log.Println("UPGRADE:", remoteIP, err)
		return
	}
	if conn, err := ws.upgrader.Upgrade(w, r, nil); err == nil {
		if ws.opts.readLimit > 0 {
			conn.SetReadLimit(ws.opts.readLimit)
		}
		s := &socket{conn: conn, sendC: make(chan []byte, ws.opts.sendCap), remoteIP: remoteIP}
		go func() {
			ws.AddSocket(s, id)
			s.read(ws)
		}()
		go s.write()
	} else {
		log.Println("UPGRADE:", remoteIP, err)
	}
}

func (ws *wsServer) Serve(addr string) error {
	s := &http.Server{Addr: addr, Handler: ws, ReadHeaderTimeout: ws.opts.readHeaderTimeout}
	if ws.opts.cert != "" && ws.opts.key != "" {
		fmt.Println("Serve WSS on", addr)
		return s.ListenAndServeTLS(ws.opts.cert, ws.opts.key)
	}
	fmt.Println("Serve WS on", addr)
	return s.ListenAndServe()
}

func NewServer(a auth.IAuth, h server.IHandler, opts ...Option) server.IServer {
	os := &options{
		readLimit:         0xFFFF,
		sendCap:           128,
		timeout:           5 * time.Minute,
		readHeaderTimeout: 5 * time.Second,
		checkOrigin:       func(r *http.Request) bool { return true },
	}
	for _, opt := range opts {
		opt(os)
	}
	if os.ipParser == nil {
		os.ipParser, _ = ip.NewParser(nil)
	}
	return &wsServer{*server.New(a, h), &websocket.Upgrader{CheckOrigin: os.checkOrigin}, os}
}
