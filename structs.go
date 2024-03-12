package engarde

import (
	"net"
)

type sendingRoutine struct {
	SrcSock   *net.UDPConn
	SrcAddr   string
	DstAddr   *net.UDPAddr
	LastRec   int64
	IsClosing bool
}

// ConnectedClient contains the information about a client
type ConnectedClient struct {
	Addr *net.UDPAddr
	Last int64
}

type UDPSocketWriter struct {
	SrcSock *net.UDPConn
	Addr    *net.UDPAddr
}

func (u *UDPSocketWriter) Write(p []byte) (n int, err error) {
	return u.SrcSock.WriteToUDP(p, u.Addr)
}

func NewUDPSocketWriter(sock *net.UDPConn, addr *net.UDPAddr) *UDPSocketWriter {
	return &UDPSocketWriter{
		SrcSock: sock,
		Addr:    addr,
	}
}
