package ip

import (
	"net"
	"net/http"
	"strings"
)

type Parser struct {
	nets []*net.IPNet
}

func (p *Parser) Parse(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	if p.isTrust(remoteIP) {
		xff := r.Header.Get("X-Forwarded-For")
		if xff == "" {
			if xri := strings.TrimSpace(r.Header.Get("X-Real-Ip")); xri != "" {
				return xri
			}
		} else {
			ips := strings.Split(xff, ",")
			for i := len(ips) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(ips[i])
				if !p.isTrust(ip) {
					return ip
				}
			}
		}
	}
	return remoteIP
}

func (p *Parser) isTrust(remoteIP string) bool {
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}
	for _, n := range p.nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func NewParser(cidrs []string) (*Parser, error) {
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		nets = append(nets, n)
	}
	return &Parser{nets}, nil
}
