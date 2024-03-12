package main

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
