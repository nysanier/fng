package ddnsimpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgconf"
	"github.com/nysanier/fng/src/pkg/pkgddns"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkglog"
)

func getDnsRR() string {
	switch pkgenv.GetEnv() {
	case pkgenv.FnEnv_Dev:
		return "dev"
	case pkgenv.FnEnv_Daily:
		return "daily"
	case pkgenv.FnEnv_Stg:
		return "stg"
	default:
		return "test"
	}
}

func UpdateDns() error {
	// 不管成功或者失败，都要求执行这个sleep
	defer func() {
		// 默认10分钟执行一次
		dnsUpdateInterval := pkgconf.GetConfigIntegerWithDefault("base", "common", "dns_update_interval", 600)
		if dnsUpdateInterval < 30 { // 至少30秒钟才执行一次
			dnsUpdateInterval = 30
		}
		pkglog.Infov("EvtDnsDumpUpdateInterval",
			"Interval", dnsUpdateInterval)
		time.Sleep(time.Second * time.Duration(dnsUpdateInterval))
	}()

	switch pkgenv.GetEnv() {
	case pkgenv.FnEnv_Dev, pkgenv.FnEnv_Daily:
	default: // 其他环境不需要自动更新
		return nil
	}

	serviceIP, err := parseServiceIP()
	if err != nil {
		pkglog.Warnv("EvtDnsParseServiceIPFail",
			"Error", err)
		return err
	}

	// 公网 ip 没有变化，因此不需要更新dns
	if serviceIP == pkgddns.CurrentServiceIP {
		return nil
	}
	pkgddns.CurrentServiceIP = serviceIP

	rr := getDnsRR()
	if err := pkgclient.SetA3927Dns(rr, serviceIP); err != nil {
		pkglog.Warnv("EvtDnsSetA3927DnsFail",
			"Error", err)
		return err
	}

	pkglog.Infov("EvtDnsUpdateDnsOK",
		"ServiceIP", serviceIP)
	return nil
}

/*
parseServiceIP 通过curl cip.cc来解析当前服务的公网ip

curl http://cip.cc
IP	: xx.xx.xx.xx
地址	: 中国  浙江  杭州
运营商	: 电信
数据二	: 浙江省杭州市 | 电信
数据三	: 中国浙江省杭州市 | 电信
URL	: http://www.cip.cc/xx.xx.xx.xx
*/
func parseServiceIP() (string, error) {
	var resp *http.Response
	var err error

	client := &http.Client{}
	req := &http.Request{
		Header: map[string][]string{
			"User-Agent": {"curl/7.79.1"},
		},
		URL: &url.URL{
			Scheme: "http",
			Host:   "cip.cc",
		},
	}

	for i := 0; i < 1; i++ {
		//resp, err = http.Get("http://cip.cc")
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		pkglog.Warnv("EvtHttpGetFail",
			"Error", err)
		//time.Sleep(time.Second * 3)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pkglog.Warnv("EvtDnsReadHttpBodyFail",
			"Error", err)
		return "", err
	}

	str := string(buf)
	lines := strings.Split(str, "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("split fail, str=%v", str)
	}

	line0 := lines[0]
	l := strings.Split(line0, ":")
	if len(l) < 2 {
		return "", fmt.Errorf("split fail, line0=%v", line0)
	}

	serviceIP := strings.TrimSpace(l[1])
	return serviceIP, nil
}
