package udp

import "net"

func ServeEcho(addr string, bufferSize uint16, echo func([]byte) []byte) error {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	buffer := make([]byte, bufferSize)
	for {
		n, raddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		if echo == nil {
			conn.WriteToUDP(buffer[:n], raddr)
		} else {
			conn.WriteToUDP(echo(buffer[:n]), raddr)
		}
	}
}
