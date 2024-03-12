package engarde

import (
	"net"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func createSendThread(ifname, sourceAddr string, wgSock *net.UDPConn, wgAddr **net.UDPAddr) {
	dst := getDstByIfname(ifname)
	dstAddr, err := net.ResolveUDPAddr("udp4", dst)
	if err != nil {
		log.Error("Can't resolve destination address '" + dst + "' for interface '" + ifname + "', not using it")
		return
	}
	srcAddr, err := net.ResolveUDPAddr("udp4", sourceAddr+":0")
	if err != nil {
		log.Error("Can't resolve source address '" + sourceAddr + "' for interface '" + ifname + "', not using it")
		return
	}
	sock, err := udpConn(srcAddr, ifname)
	if err != nil {
		log.Error("Can't create socket for address '" + sourceAddr + "' on interface '" + ifname + "', not using it")
		return
	}

	routine := sendingRoutine{
		SrcSock:   sock,
		SrcAddr:   sourceAddr,
		DstAddr:   dstAddr,
		IsClosing: false,
	}
	ptrRoutine := &routine

	go wgWriteBack(ifname, ptrRoutine, wgSock, wgAddr)
	sendingChannelsMutex.Lock()
	sendingChannels[ifname] = ptrRoutine
	sendingChannelsMutex.Unlock()
}

func receiveFromClient(socket, wgSocket *net.UDPConn, wgAddr *net.UDPAddr) {
	buffer := make([]byte, parsedConfig.Server.MTU)
	var currentTime int64
	var n int
	var srcAddr *net.UDPAddr
	var srcAddrS string
	var client *ConnectedClient
	var exists bool
	var err error
	for {
		n, srcAddr, err = socket.ReadFromUDP(buffer)
		if err != nil {
			log.Warn("Error reading from client")
			continue
		}

		// Check if client exists
		currentTime = time.Now().Unix()
		srcAddrS = srcAddr.IP.String() + ":" + strconv.Itoa(srcAddr.Port)
		clientsMutex.RLock()
		client, exists = clients[srcAddrS]
		clientsMutex.RUnlock()
		if exists {
			client.Last = currentTime
		} else {
			log.Info("New client connected: '" + srcAddrS + "'")
			newClient := ConnectedClient{
				Addr: srcAddr,
				Last: currentTime,
			}
			clientsMutex.Lock()
			clients[srcAddrS] = &newClient
			clientsMutex.Unlock()
		}
		
		_, err = wgSocket.WriteToUDP(buffer[:n], wgAddr)
		if err != nil {
			log.Warn("Error writing to WireGuard")
		}
	}
}
