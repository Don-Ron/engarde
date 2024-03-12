package engarde

import (
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func wgWriteBack(ifname string, routine *sendingRoutine, wgSock *net.UDPConn, wgAddr **net.UDPAddr) {
	buffer := make([]byte, parsedConfig.Client.MTU)
	var n int
	var err error

	w := NewUDPSocketWriter(wgSock, *wgAddr)
	tr := io.TeeReader(routine.SrcSock, w)

	for {
		if parsedConfig.Client.UseTeeReader {
			_, err = tr.Read(buffer)
			if err != nil {
				log.Warn("Error reading from '" + ifname + "', re-creating socket")
				terminateRoutine(routine, ifname, true)
				return
			}
		} else {
			n, _, err = routine.SrcSock.ReadFromUDP(buffer)
			if err != nil {
				log.Warn("Error reading from '" + ifname + "', re-creating socket")
				terminateRoutine(routine, ifname, true)
				return
			}

			_, err = wgSock.WriteToUDP(buffer[:n], *wgAddr)
			if err != nil {
				log.Warn("Error writing to WireGuard")
			}
		}

		if routine.IsClosing {
			return
		}

		routine.LastRec = time.Now().Unix()
	}
}

func receiveFromWireguardClient(wgsock *net.UDPConn, sourceAddr **net.UDPAddr) {
	buffer := make([]byte, parsedConfig.Client.MTU)
	var n int
	var srcAddr *net.UDPAddr
	var routine *sendingRoutine
	var err error
	var ifname string
	var toDelete []string
	for {
		n, srcAddr, err = wgsock.ReadFromUDP(buffer)
		if err != nil {
			log.Warn("Error reading from Wireguard")
			continue
		}
		*sourceAddr = srcAddr
		sendingChannelsMutex.RLock()
		for ifname, routine = range sendingChannels {
			if parsedConfig.Client.WriteTimeout > 0 {
				err = routine.SrcSock.SetWriteDeadline(time.Now().Add(parsedConfig.Client.WriteTimeout * time.Millisecond))
				if err != nil {
					log.WithError(err).Warn("Error setting source socket write deadline to " + parsedConfig.Client.WriteTimeout.String())
				}
			}
			_, err = routine.SrcSock.WriteToUDP(buffer[:n], routine.DstAddr)
			if err != nil {
				log.Warn("Error writing to '" + ifname + "', re-creating socket")
				terminateRoutine(routine, ifname, false)
				toDelete = append(toDelete, ifname)
			}
		}
		sendingChannelsMutex.RUnlock()
		sendingChannelsMutex.Lock()
		for _, ifname = range toDelete {
			delete(sendingChannels, ifname)
		}
		toDelete = toDelete[:0]
		sendingChannelsMutex.Unlock()
	}
}

func receiveFromEngardeClient(wgSocket, socket *net.UDPConn) {
	buffer := make([]byte, parsedConfig.Server.MTU)
	var n int
	var client *ConnectedClient
	var currentTime int64
	var clientAddr string
	var err error
	var toDelete []string
	for {
		n, _, err = wgSocket.ReadFromUDP(buffer)
		if err != nil {
			log.Warn("Error reading from WireGuard")
			continue
		}
		currentTime = time.Now().Unix()
		clientsMutex.RLock()
		for clientAddr, client = range clients {
			if client.Last > currentTime-parsedConfig.Server.ClientTimeout {
				if parsedConfig.Server.WriteTimeout > 0 {
					err = socket.SetWriteDeadline(time.Now().Add(parsedConfig.Server.WriteTimeout * time.Millisecond))
					if err != nil {
						log.WithError(err).Warn("Error setting write deadline to " + parsedConfig.Server.WriteTimeout.String())
					}
				}
				_, err = socket.WriteToUDP(buffer[:n], client.Addr)
				if err != nil {
					log.Warn("Error writing to client '" + clientAddr + "', terminating it")
					toDelete = append(toDelete, clientAddr)
				}
			} else {
				log.Info("Client '" + clientAddr + "' timed out")
				toDelete = append(toDelete, clientAddr)
			}
		}
		clientsMutex.RUnlock()
		clientsMutex.Lock()
		for _, clientAddr = range toDelete {
			delete(clients, clientAddr)
		}
		clientsMutex.Unlock()
		toDelete = toDelete[:0]
	}
}
