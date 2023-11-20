package dns

import "errors"

type ManageDNS interface {
	GetIpV4(domain string) (string, error)
	GetIpV6(domain string) (string, error)
	UpdateIpV4(domain string, ipv4 string) error
	UpdateIpV6(domain string, ipv6 string) error
}

func GetDnsManager(t string) (ManageDNS, error) {
	switch t {
	case "aliyun":
		return NewAliyunDns()
	default:
		return nil, errors.New("not support ip detector: " + t)
	}
}
