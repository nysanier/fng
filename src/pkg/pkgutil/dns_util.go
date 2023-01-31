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
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

var (
	ServiceIP = "11.22.33.44" // 定时更新
)

const (
	FormatTime = "15:4:5"
)

func GetServiceIP() string {
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
		log.Printf("ServiceIP is %v", ServiceIP)
	}

	// else use the real service ip
	return ServiceIP
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

// TODO: 更新之前可以先比较一下，不过当前更新的频率比较低，性能也是ok的
func updateDns() error {
	switch pkgvar.FnEnv {
	case pkgvar.FnEnv_Dev, pkgvar.FnEnv_Daily:
	default: // 其他环境不需要自动更新
		return nil
	}

	if err := parseServiceIP(); err != nil {
		log.Printf("parseServiceIP fail, err=%v", err)
		return err
	}

	rr := getDnsRR()
	ip := GetServiceIP()
	if err := pkgclient.SetA3927Dns(rr, ip); err != nil {
		log.Printf("SetA3927Dns fail, err=%v", err)
		return err
	}

	return nil
}

func RunDnsUpdater() {
	ticker := time.NewTicker(time.Hour)

	// 启动的时候先执行一次
	if err := updateDns(); err != nil {
		log.Printf("updateDns fail, err=%v", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			//log.Printf("test ticker")
			if err := updateDns(); err != nil {
				log.Printf("updateDns fail, err=%v", err)
				time.Sleep(time.Second * 3)
				continue
			}
		}
	}

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
func parseServiceIP() error {
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

	for i := 0; i < 3; i++ {
		//resp, err = http.Get("http://cip.cc")
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		log.Printf("http.Get fail, err: %v", err)
		time.Sleep(time.Second * 3)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll fail, err: %v", err)
		return err
	}

	str := string(buf)
	lines := strings.Split(str, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("split fail, str=%v", str)
	}

	line0 := lines[0]
	l := strings.Split(line0, ":")
	if len(l) < 2 {
		return fmt.Errorf("split fail, line0=%v", line0)
	}

	ServiceIP = strings.TrimSpace(l[1])
	return nil
}
