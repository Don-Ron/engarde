package engarde

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

var sendingChannels = make(map[string]*sendingRoutine)
var sendingChannelsMutex = &sync.RWMutex{}

func RunClient(configName string) {
	parseConfig(Client, configName)

	var wireguardAddr *net.UDPAddr
	ptrWireguardAddr := &wireguardAddr

	WireguardListenAddr, err := net.ResolveUDPAddr("udp4", parsedConfig.Client.ListenAddr)
	handleErr(err, "error resolving listen address for wireguard")
	WireguardSocket, err := net.ListenUDP("udp", WireguardListenAddr)
	handleErr(err, "error creating wireguard socket")
	log.Info("Listening on " + parsedConfig.Client.ListenAddr)

	go updateAvailableInterfaces(WireguardSocket, ptrWireguardAddr)

	receiveFromWireguardClient(WireguardSocket, ptrWireguardAddr)
}
