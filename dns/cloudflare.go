package dns

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareDns struct {
	api *cloudflare.API

	zoneIDmap map[string]string
	rrMap     map[string]string
}

func NewCloudflareDns() (*CloudflareDns, error) {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if apiToken == "" {
		return nil, errors.New("need env CLOUDFLARE_API_TOKEN")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, err
	}

	return &CloudflareDns{
		api:       api,
		zoneIDmap: make(map[string]string),
		rrMap:     make(map[string]string),
	}, nil
}

func (d *CloudflareDns) GetIpV4(domain, rr string) (string, error) {
	record, err := d.getDnsRecord(domain, rr, "A")
	if err != nil {
		return "", err
	}
	return record.Content, nil
}

func (d *CloudflareDns) GetIpV6(domain, rr string) (string, error) {
	record, err := d.getDnsRecord(domain, rr, "AAAA")
	if err != nil {
		return "", err
	}
	return record.Content, nil
}

func (d *CloudflareDns) UpdateIpV4(domain, rr string, ipv4 string) error {
	return d.updateIP(domain, rr, "A", ipv4)
}

func (d *CloudflareDns) UpdateIpV6(domain, rr string, ipv6 string) error {
	return d.updateIP(domain, rr, "AAAA", ipv6)
}

func (d *CloudflareDns) getZoneID(domain string) (string, error) {
	if zoneID, ok := d.zoneIDmap[domain]; ok {
		return zoneID, nil
	}
	zoneID, err := d.api.ZoneIDByName(domain)
	if err != nil {
		return "", err
	}
	d.zoneIDmap[domain] = zoneID
	return zoneID, nil
}

func (d *CloudflareDns) updateIP(domain, rr, recordType, ip string) error {
	zoneID, err := d.getZoneID(domain)
	if err != nil {
		return err
	}

	record, err := d.getDnsRecord(domain, rr, recordType)
	if err != nil {
		return err
	}

	_, err = d.api.UpdateDNSRecord(context.Background(), cloudflare.ResourceIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
		ID:      record.ID,
		Content: ip,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *CloudflareDns) toName(domain, rr string) string {
	if rr == "@" || rr == domain || rr == "" {
		return domain
	}
	return fmt.Sprintf("%s.%s", rr, domain)
}

func (d *CloudflareDns) getDnsRecord(domain, rr, recordType string) (cloudflare.DNSRecord, error) {
	zoneID, err := d.getZoneID(domain)
	if err != nil {
		return cloudflare.DNSRecord{}, err
	}
	records, _, err := d.api.ListDNSRecords(
		context.Background(),
		cloudflare.ResourceIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{
			Name: d.toName(domain, rr),
			Type: recordType,
		},
	)
	if err != nil {
		return cloudflare.DNSRecord{}, err
	}

	if len(records) == 0 {
		return cloudflare.DNSRecord{}, errors.New("no record found")
	}

	if len(records) > 1 {
		return cloudflare.DNSRecord{}, errors.New("more than one record found")
	}

	return records[0], nil

}
