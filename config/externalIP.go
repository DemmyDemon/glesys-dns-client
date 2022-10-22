package config

import (
	"errors"
	"net"
	"strings"
)

var ErrCouldNotDetermine = errors.New("could not detemine your IPv4 address")

// GetExternalIPv4 returns the first found external IPv4 of the current host.
// "external" meaning it is not a 192.168.0.0/16 or a 10.0.0.0/24
func GetExternalIPv4() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue // We don't care about loopback
		}
		if iface.Flags&net.FlagUp == 0 {
			continue // We only care if it's up
		}
		addresses, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, address := range addresses {
			var ip net.IP
			switch value := address.(type) {
			case *net.IPNet:
				ip = value.IP
			case *net.IPAddr:
				ip = value.IP
			}
			if ip == nil {
				continue // Irrelevant interface type
			}
			if ip.IsLoopback() {
				continue // Yeah, we already checked this for the IF, but yaknow, belt *and* suspenders
			}
			ip = ip.To4()
			if ip == nil {
				continue // We only want IPv4!
			}
			ipString := ip.String()
			if strings.HasPrefix(ipString, "192.168.") || strings.HasPrefix(ipString, "10.") {
				continue // Local network addresses don't count!
			}
			return ipString, nil
		}
	}
	return "", ErrCouldNotDetermine
}
