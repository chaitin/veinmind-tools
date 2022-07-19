package netutil

import (
	"github.com/shirou/gopsutil/net"
	"regexp"
	"strconv"
)

func GetPidByPort(port string) (int, error) {
	netConnections, err := net.Connections("all")
	if err != nil {
		return 0, err
	}
	for _, connection := range netConnections {
		if strconv.Itoa(int(connection.Laddr.Port)) == port {
			return int(connection.Pid), nil
		}
	}
	return 0, nil
}

func IsLocalHost(host string) (bool, error) {
	localIPs, err := GetLocalIP()
	if err != nil {
		return false, err
	}
	for _, localIP := range localIPs {
		if host == localIP {
			return true, nil
		}
	}
	return false, nil
}

// GetLocalIP get all loacl IP
func GetLocalIP() ([]string, error) {
	addresses, err := net.Interfaces()
	reg := `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
	matchIP := regexp.MustCompile(reg)
	if err != nil {
		return nil, err
	}
	IPs := make([]string, 0)
	for _, address := range addresses {
		for _, addr := range address.Addrs {
			if matchIP.MatchString(addr.Addr) {
				IPs = append(IPs, matchIP.FindString(addr.Addr))
			}
		}
	}
	return IPs, nil
}
