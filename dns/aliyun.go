package dns

import (
	"errors"
	"log"
	"os"

	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
)

type AliyunDns struct {
	client *dns.Client
	// domain+rr+iptype => recordId map
	// 减少更新时的api请求次数
	rrMap map[string]string
}

func getClient(regionId *string, accessKeyId *string, accessKeySecret *string) (*dns.Client, error) {
	config := &openapi.Config{}
	// 您的AccessKey ID
	config.AccessKeyId = accessKeyId
	// 您的AccessKey Secret
	config.AccessKeySecret = accessKeySecret

	// 您的可用区ID
	config.RegionId = regionId

	_result, _err := dns.NewClient(config)
	if _err != nil {
		return nil, _err
	}
	return _result, nil
}
func getTargetDNSRecord(client *dns.Client, domainName *string, RR *string, recordType *string) (*dns.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	req := &dns.DescribeDomainRecordsRequest{}
	// 主域名
	req.DomainName = domainName
	// 主机记录
	req.RRKeyWord = RR
	// 解析记录类型
	req.Type = recordType

	resp, _err := client.DescribeDomainRecords(req)
	if _err != nil {
		return nil, _err
	}

	log.Println("aliyun: ", resp.Body)

	if resp.Body.DomainRecords.Record != nil && len(resp.Body.DomainRecords.Record) != 1 {
		return nil, errors.New("records length!= 1")
	}
	return resp.Body.DomainRecords.Record[0], nil
}

//   req := &dns.UpdateDomainRecordRequest{}
//   // 主机记录
//   req.RR = RR
//   // 记录ID
//   req.RecordId = recordId
//   // 将主机记录值改为当前主机IP
//   req.Value = currentHostIP
//   // 解析记录类型
//   req.Type = recordType
/**
 * 修改解析记录
 */
func updateDomainRecord(client *dns.Client, recordId string, rr string, recordType string, value string) error {
	req := &dns.UpdateDomainRecordRequest{}
	req.RecordId = &recordId
	req.RR = &rr
	req.Value = &value
	req.Type = &recordType

	resp, _err := client.UpdateDomainRecord(req)
	if _err != nil {
		return _err
	}

	log.Println("aliyun: ", resp.Body)
	return nil
}

func NewAliyunDns() (*AliyunDns, error) {
	// accessKeyId, accessKeySecret, domain, regionId string
	accessKeyId := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	regionId := os.Getenv("ALIYUN_REGION_ID")
	if accessKeyId == "" {
		return nil, errors.New("need env ALIYUN_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		return nil, errors.New("need env ALIYUN_ACCESS_KEY_SECRET")
	}
	if regionId == "" {
		return nil, errors.New("need env ALIYUN_REGION_ID")
	}

	client, err := getClient(&regionId, &accessKeyId, &accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliyunDns{
		client: client,
		rrMap:  make(map[string]string),
	}, nil
}

func (d *AliyunDns) UpdateIpV4(domain, rr, ipv4 string) error {

	rcdtype := "A"
	recordidmapkey := domain + rr + rcdtype
	_, ok := d.rrMap[recordidmapkey]
	if !ok {
		_, err := d.GetIpV4(domain, rr)
		if err != nil {
			return err
		}
	}
	recordid, ok := d.rrMap[recordidmapkey]
	if !ok {
		return errors.New("get recordid failed")
	}
	if err := updateDomainRecord(d.client, recordid, rr, rcdtype, ipv4); err != nil {
		return err
	}

	return nil
}

func (d *AliyunDns) UpdateIpV6(domain, rr, ipv6 string) error {

	rcdtype := "AAAA"
	recordidmapkey := domain + rr + rcdtype
	_, ok := d.rrMap[recordidmapkey]
	if !ok {
		_, err := d.GetIpV4(domain, rr)
		if err != nil {
			return err
		}
	}
	recordid, ok := d.rrMap[recordidmapkey]
	if !ok {
		return errors.New("get recordid failed")
	}
	if err := updateDomainRecord(d.client, recordid, rr, rcdtype, ipv6); err != nil {
		return err
	}

	return nil
}

func (d *AliyunDns) GetIpV4(domain, rr string) (string, error) {

	rcdtype := "A"
	record, err := getTargetDNSRecord(d.client, &domain, &rr, &rcdtype)

	if err != nil {
		return "", err
	}

	d.rrMap[domain+rr+rcdtype] = *record.RecordId

	return *record.Value, nil
}

func (d *AliyunDns) GetIpV6(domain, rr string) (string, error) {

	rcdtype := "AAAA"
	record, err := getTargetDNSRecord(d.client, &domain, &rr, &rcdtype)
	if err != nil {
		return "", err
	}

	d.rrMap[domain+rr+rcdtype] = *record.RecordId

	return *record.Value, nil
}
