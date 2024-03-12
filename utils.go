package engarde

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

var exclusionSwaps = make(map[string]bool)

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s | %s", msg, err))
	}
}

func isSwapped(name string) bool {
	if _, ok := exclusionSwaps[name]; ok {
		return true
	}
	return false
}

func isExcluded(name string) bool {
	if len(parsedConfig.Client.IncludedInterfaces) > 0 {
		for _, ifname := range parsedConfig.Client.IncludedInterfaces {
			if ifname == name {
				return false
			}
		}

		return true
	}

	for _, ifname := range parsedConfig.Client.ExcludedInterfaces {
		if ifname == name {
			return !isSwapped(name)
		}
	}

	return isSwapped(name)
}

func isAddressAllowed(addr string) bool {
	// TODO: IPv6 support
	if strings.ContainsRune(addr, ':') {
		return false
	}
	ip := net.ParseIP(addr)
	disallowedNetworks := []string{
		"169.254.0.0/16",
		"127.0.0.0/8",
	}
	for _, disallowedNetwork := range disallowedNetworks {
		_, subnet, _ := net.ParseCIDR(disallowedNetwork)
		if subnet.Contains(ip) {
			return false
		}
	}
	return true
}

func terminateRoutine(routine *sendingRoutine, ifname string, deleteFromSlice bool) {
	routine.IsClosing = true
	routine.SrcSock.Close()
	if deleteFromSlice {
		sendingChannelsMutex.Lock()
		delete(sendingChannels, ifname)
		sendingChannelsMutex.Unlock()
	}
}
