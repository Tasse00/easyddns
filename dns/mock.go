package dns

type DnsManageMock struct {
}

func NewDnsManageMock() *DnsManageMock {
	return &DnsManageMock{}
}

func (d *DnsManageMock) GetIpV4(domain, rr string) (string, error) {
	return "127.0.0.1", nil
}

func (d *DnsManageMock) GetIpV6(domain, rr string) (string, error) {
	return "::1", nil
}

func (d *DnsManageMock) UpdateIpV4(domain, rr, ipv4 string) error {
	return nil
}

func (d *DnsManageMock) UpdateIpV6(domain, rr, ipv6 string) error {
	return nil
}
