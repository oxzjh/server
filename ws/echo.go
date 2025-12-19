package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type echoServer struct {
	upgrader *websocket.Upgrader
	echo     func([]byte) []byte
}

func (es *echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conn, err := es.upgrader.Upgrade(w, r, nil); err == nil {
		go es.echoWS(conn)
	}
}

func (es *echoServer) echoWS(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if es.echo == nil {
			conn.WriteMessage(websocket.BinaryMessage, data)
		} else {
			conn.WriteMessage(websocket.BinaryMessage, es.echo(data))
		}
	}
}

func ServeEcho(addr string, echo func([]byte) []byte) error {
	return http.ListenAndServe(addr, &echoServer{&websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}, echo})
}
