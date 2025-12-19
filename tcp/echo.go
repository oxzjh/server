package tcp

import "net"

func ServeEcho(addr string, bufferSize int, echo func([]byte) []byte) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go echoTCP(conn, bufferSize, echo)
	}
}

func echoTCP(conn net.Conn, bufferSize int, echo func([]byte) []byte) {
	buffer := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		if echo == nil {
			conn.Write(buffer[:n])
		} else {
			conn.Write(echo(buffer[:n]))
		}
	}
}
