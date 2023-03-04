package pkgutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgconf"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

var (
	CurrentServiceIP = "11.22.33.44" // 定时更新
)

const (
	FormatTime = "15:4:5"
)

func GetCurrentServiceIP() string {
	if pkgvar.IsDevEnv() {
		// 将当前时间格式化为ip，方便观察
		//s := pkgfunc.GetCstNow().Format(FormatTime)
		//
		//l := strings.Split(s, ":")
		//
		//if len(l) < 3 {
		//	return "0.0.0.0"
		//}
		//
		//hour, _ := strconv.ParseInt(l[0], 10, 64)
		//min, _ := strconv.ParseInt(l[1], 10, 64)
		//sec, _ := strconv.ParseInt(l[2], 10, 64)
		//return fmt.Sprintf("0.%v.%v.%v", hour, min, sec)
		//log.Printf("CurrentServiceIP is %v", CurrentServiceIP)
	}

	// else use the real service ip
	return CurrentServiceIP
}

func getDnsRR() string {
	switch pkgvar.FnEnv {
	case pkgvar.FnEnv_Dev:
		return "dev"
	case pkgvar.FnEnv_Daily:
		return "daily"
	case pkgvar.FnEnv_Stg:
		return "stg"
	default:
		return "test"
	}
}

func updateDns() error {
	// 不管成功或者失败，都要求执行这个sleep
	defer func() {
		// 默认10分钟执行一次
		dnsUpdateInterval := pkgconf.GetConfigIntegerWithDefault("base", "common", "dns_update_interval", 600)
		if dnsUpdateInterval < 30 { // 至少30秒钟才执行一次
			dnsUpdateInterval = 30
		}
		log.Printf("dnsUpdateInterval: %v", dnsUpdateInterval)
		time.Sleep(time.Second * time.Duration(dnsUpdateInterval))
	}()

	switch pkgvar.FnEnv {
	case pkgvar.FnEnv_Dev, pkgvar.FnEnv_Daily:
	default: // 其他环境不需要自动更新
		return nil
	}

	serviceIP, err := parseServiceIP()
	if err != nil {
		log.Printf("parseServiceIP fail, err=%v", err)
		return err
	}

	// 公网 ip 没有变化，因此不需要更新dns
	if serviceIP == CurrentServiceIP {
		return nil
	}
	CurrentServiceIP = serviceIP

	rr := getDnsRR()
	if err := pkgclient.SetA3927Dns(rr, serviceIP); err != nil {
		log.Printf("SetA3927Dns fail, err=%v", err)
		return err
	}

	log.Printf("updateDns ok, serviceIP=%v", serviceIP)
	return nil
}

var (
	dnsUpdateTimer *pkgfunc.Timer
)

func StartDnsUpdater() {
	dnsUpdateTimer = pkgfunc.NewTimer(updateDns, time.Duration(0)) // 由updateDns来控制时间间隔
	dnsUpdateTimer.SetFirstDelay(time.Second * 5)
	dnsUpdateTimer.Start()
	log.Printf("start dns updater ok")
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

		log.Printf("http.Get fail, err: %v", err)
		//time.Sleep(time.Second * 3)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll fail, err: %v", err)
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
