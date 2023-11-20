package ip

import (
	"errors"
	"net"
)

type DeviceInterface struct{}

func NewDeviceInterface() *DeviceInterface {
	return &DeviceInterface{}
}

// use first notloopback ipv4
func (d *DeviceInterface) DetectIpV4() (string, error) {

	targetV := "v4"
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				v := ""
				if ipNet.IP.To4() == nil && ipNet.IP.IsGlobalUnicast() {
					v = "v6"
				} else if ipNet.IP.To4() != nil {
					v = "v4"
				}

				if targetV == v {
					return ipNet.IP.String(), nil
				}
			}

		}
	}

	return "", errors.New("not ip addr found")
}

// use first notloopback ipv6
func (d *DeviceInterface) DetectIpV6() (string, error) {

	targetV := "v6"
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				v := ""
				if ipNet.IP.To4() == nil && ipNet.IP.IsGlobalUnicast() {
					v = "v6"
				} else if ipNet.IP.To4() != nil {
					v = "v4"
				}

				if targetV == v {
					return ipNet.IP.String(), nil
				}
			}

		}
	}

	return "", errors.New("not ip addr found")
}
