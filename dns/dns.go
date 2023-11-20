package dns

import "errors"

type ManageDNS interface {
	GetIpV4(domain, rr string) (string, error)
	GetIpV6(domain, rr string) (string, error)
	UpdateIpV4(domain, rr string, ipv4 string) error
	UpdateIpV6(domain, rr string, ipv6 string) error
}

func GetDnsManager(t string) (ManageDNS, error) {
	switch t {
	case "aliyun":
		return NewAliyunDns()
	default:
		return nil, errors.New("not support ip detector: " + t)
	}
}
