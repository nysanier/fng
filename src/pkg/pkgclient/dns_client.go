package pkgclient

import (
	"log"
	"strings"
	"sync"

	dnsclient "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openclient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/nysanier/fng/src/pkg/pkgconfig"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
)

var (
	dnsClientOnce sync.Once
	dnsClient     *dnsclient.Client
)

func GetDnsClient() *dnsclient.Client {
	cfg := &openclient.Config{
		AccessKeyId:     pkgfunc.StringPtr(pkgconfig.AK),
		AccessKeySecret: pkgfunc.StringPtr(pkgconfig.SK),
		RegionId:        pkgfunc.StringPtr("cn-hangzhou"),
	}

	dnsClientOnce.Do(func() {
		var err error
		dnsClient, err = dnsclient.NewClient(cfg)
		if err != nil {
			log.Printf("dnsclient.NewClient fail, err=%v", err)
			panic(err)
		}
	})

	return dnsClient
}

func ListA3927Dns() string {
	client := GetDnsClient()
	req := &dnsclient.DescribeDomainRecordsRequest{
		DomainName: pkgfunc.StringPtr("a3927.top"),
	}
	resp, err := client.DescribeDomainRecords(req)
	if err != nil {
		log.Printf("client.DescribeDomainRecords fail, err=%v", err)
	}
	for i, v := range resp.Body.DomainRecords.Record {
		log.Printf("record: %v <-> %v", i, pkgfunc.FormatJson(v))
	}

	return "TODO"
}

var (
	A3927DnsMap = map[string]string{
		"dev":   "808500502808423424",
		"daily": "806418358604756992",
		"stg":   "807758159302026240",
	}
)

// 设置到a3927.top这个domain
func SetA3927Dns(rr, value string) error {
	recordID := A3927DnsMap[rr]

	client := GetDnsClient()
	req := &dnsclient.UpdateDomainRecordRequest{
		RR:       pkgfunc.StringPtr(rr),
		RecordId: pkgfunc.StringPtr(recordID),
		Type:     pkgfunc.StringPtr("A"),
		Value:    pkgfunc.StringPtr(value),
	}
	resp, err := client.UpdateDomainRecord(req)
	if err != nil {
		// 已经存在重复的记录，因此认为执行成功
		if strings.Contains(err.Error(), "DomainRecordDuplicate") {
			return nil
		}

		log.Printf("client.UpdateDomainRecord fail, err=%v", err)
		return err
	}

	_ = resp
	return nil
}
