package system

import (
	"fmt"
	nnet "net"
	"time"

	"github.com/shirou/gopsutil/net"
)

// OutboundIP get local PC outbound IP.
func OutboundIP(server string) (string, error) {
	conn, err := nnet.DialTimeout("udp", server, time.Duration(1500)*time.Millisecond)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*nnet.UDPAddr)

	return localAddr.IP.String(), nil
}

// LocalIP get local PC IP.
func LocalIP() (string, error) {
	addrs, err := nnet.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipAddr, ok := addr.(*nnet.IPNet)
		if !ok {
			continue
		}

		fmt.Printf("ip = %s\n", ipAddr.IP.String())

		if ipAddr.IP.IsLoopback() {
			continue
		}

		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}

		return ipAddr.IP.String(), nil
	}

	return "", fmt.Errorf("not found")
}

// NetInfo get net interface IO info.
func NetInfo() ([]net.IOCountersStat, error) {
	return net.IOCounters(true)
}
