package ip

import "errors"

type DetectIp interface {
	DetectIpV4() (string, error)
	DetectIpV6() (string, error)
}

func GetIpDetector(t string) (DetectIp, error) {
	switch t {
	case "netarm":
		return &Netarm{}, nil
	default:
		return nil, errors.New("not support ip detector: " + t)
	}
}
