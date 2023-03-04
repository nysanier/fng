package pkgclient

import (
	"strings"
	"sync"

	dnsclient "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openclient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

var (
	dnsClientOnce sync.Once
	dnsClient     *dnsclient.Client
)

func GetDnsClient() *dnsclient.Client {
	cfg := &openclient.Config{
		AccessKeyId:     pkgfunc.StringPtr(pkgvar.GetAK()),
		AccessKeySecret: pkgfunc.StringPtr(pkgvar.GetSK()),
		RegionId:        pkgfunc.StringPtr(pkgvar.RegionID),
	}

	dnsClientOnce.Do(func() {
		var err error
		dnsClient, err = dnsclient.NewClient(cfg)
		if err != nil {
			pkglog.Warnv("EvtDnsNewClientFail",
				"Error", err)
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
		pkglog.Warnv("EvtDnsDescribeDomainRecordsFail",
			"Error", err)
	}
	for i, v := range resp.Body.DomainRecords.Record {
		pkglog.Infov("EvtDnsDumpDomainRecord",
			"Index", i,
			"Record", pkgfunc.FormatJson(v))
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

		pkglog.Warnv("EvtDnsUpdateDomainRecordFail",
			"Error", err)
		return err
	}

	_ = resp
	return nil
}
