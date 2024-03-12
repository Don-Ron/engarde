package engarde

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

var clients map[string]*ConnectedClient
var clientsMutex *sync.RWMutex

func RunServer(configName string) {
	parseConfig(Server, configName)

	clients = make(map[string]*ConnectedClient)
	clientsMutex = &sync.RWMutex{}

	WireguardAddr, err := net.ResolveUDPAddr("udp4", parsedConfig.Server.DstAddr)
	handleErr(err, "cannot resolve wireguard destination address")
	WireguardSource, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	handleErr(err, "cannot resolve wireguard listen address")
	WireguardSocket, err := net.ListenUDP("udp", WireguardSource)
	handleErr(err, "cannot initialize wireguard socket")

	ClientsListenAddr, err := net.ResolveUDPAddr("udp4", parsedConfig.Server.ListenAddr)
	handleErr(err, "cannot resolve engarde listen address")
	ClientSocket, err := net.ListenUDP("udp", ClientsListenAddr)
	handleErr(err, "cannot create engarde listen socket")
	log.Info("Listening on " + parsedConfig.Server.ListenAddr)

	go receiveFromEngardeClient(WireguardSocket, ClientSocket)
	receiveFromClient(ClientSocket, WireguardSocket, WireguardAddr)
}
