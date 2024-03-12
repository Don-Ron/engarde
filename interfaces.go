package engarde

import (
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func interfaceExists(interfaces []net.Interface, name string) bool {
	for _, iface := range interfaces {
		if iface.Name == name {
			return true
		}
	}
	return false
}

func getAddressByInterface(iface net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		splAddr := strings.Split(addr.String(), "/")[0]
		if isAddressAllowed(splAddr) {
			return splAddr
		}
	}
	return ""
}

func getDstByIfname(ifname string) string {
	for _, override := range parsedConfig.Client.DstOverrides {
		if override.IfName == ifname {
			return override.DstAddr
		}
	}
	return parsedConfig.Client.DstAddr
}

func ListInterfaces() {
	interfaces, err := net.Interfaces()
	handleErr(err, "listInterfaces 1")
	for _, iface := range interfaces {
		ifname := iface.Name
		print("\r\n" + ifname + "\r\n")
		ifaddr := getAddressByInterface(iface)
		print("  Address: " + ifaddr + "\r\n")
	}
}

func updateAvailableInterfaces(wgSock *net.UDPConn, wgAddr **net.UDPAddr) {
	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// Delete unavailable interfaces
		for ifname, routine := range sendingChannels {
			if !interfaceExists(interfaces, ifname) {
				log.Info("Interface '" + ifname + "' no longer exists, deleting it")
				terminateRoutine(routine, ifname, true)
				continue
			}
			if isExcluded(ifname) {
				log.Info("Interface '" + ifname + "' is now excluded, deleting it")
				terminateRoutine(routine, ifname, true)
				continue
			}
			iface, err := net.InterfaceByName(ifname)
			if err != nil {
				continue
			}
			ifaddr := getAddressByInterface(*iface)
			if ifaddr != routine.SrcAddr {
				log.Info("Interface '" + ifname + "' changed address, re-creating socket")
				terminateRoutine(routine, ifname, true)
				continue
			}
		}
		for _, iface := range interfaces {
			ifname := iface.Name
			if isExcluded(ifname) {
				continue
			}
			if _, ok := sendingChannels[ifname]; ok {
				continue
			}
			ifaddr := getAddressByInterface(iface)
			if ifaddr != "" {
				log.Info("New interface '" + ifname + "' with IP '" + ifaddr + "', adding it")
				createSendThread(ifname, getAddressByInterface(iface), wgSock, wgAddr)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
