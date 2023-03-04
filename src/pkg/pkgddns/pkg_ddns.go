package pkgddns

import (
	"time"

	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
)

var (
	CurrentServiceIP = "11.22.33.44" // 定时更新
)

var (
	dnsUpdateTimer *pkgfunc.Timer
)

func InitDdns(updateDns func() error) {
	dnsUpdateTimer = pkgfunc.NewTimer(updateDns, time.Duration(0)) // 由updateDns来控制时间间隔
	interval := time.Second * 5
	dnsUpdateTimer.SetFirstDelay(interval)
	dnsUpdateTimer.Start()
	pkglog.Infov("EvtDdnsInitOK",
		"UpdateInterval", interval)
}

func GetCurrentServiceIP() string {
	if pkgenv.IsDevEnv() {
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
