package main

import (
	"net/http"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgconf/confimpl"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkglog/logimpl"
	"github.com/nysanier/fng/src/pkg/pkgutil"
	"github.com/nysanier/fng/src/pkg/pkgvar"
	"github.com/nysanier/fng/src/pkg/version"
	"github.com/nysanier/fng/src/svcs/echo-svc/router"
)

func main() {
	// 初始化 随机数、时区、启动时间
	pkgfunc.InitRand()
	pkgvar.TzLoc = pkgfunc.LoadTzLoc()
	pkgvar.FnStartTime = getCRFC3339CstTimeStr()

	// 初始化 env/log
	pkgenv.LoadEnv()
	logimpl.InitSlsLog()

	pkglog.Infov("EvtFngInitBegin",
		"AppVersion", version.GetAppVersion(),
		"GitCommit", version.GetShortGitCommit(),
		"BuildTime", version.GetBuildTimeStr())

	// 初始化 conf/dns
	confimpl.StartConfigUpdater()
	pkgutil.StartDnsUpdater()

	// Start Http Server
	r := router.InitRouter()
	server := &http.Server{
		Addr:           ":17080",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 2 * 1024 * 1024,
	}
	go server.ListenAndServe()

	pkglog.Infov("EvtFngInitEnd")
	var ch chan int
	<-ch
}

func getCRFC3339CstTimeStr() string {
	t := pkgfunc.GetCstNow()
	str := pkgfunc.GetRFC3339TimeStr(t)
	return str
}
