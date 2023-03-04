package main

import (
	"net/http"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgconf"
	"github.com/nysanier/fng/src/pkg/pkgconf/confimpl"
	"github.com/nysanier/fng/src/pkg/pkgddns"
	"github.com/nysanier/fng/src/pkg/pkgddns/ddnsimpl"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkglog/logimpl"
	"github.com/nysanier/fng/src/pkg/pkgutil"
	"github.com/nysanier/fng/src/pkg/version"
	"github.com/nysanier/fng/src/svcs/echo-svc/router"
)

func main() {
	// 初始化 随机数、时区、启动时间
	pkgfunc.InitRand()
	pkgutil.InitTzLoc()
	pkgutil.InitStartTime()

	// 初始化 env/log
	pkgenv.InitEnv()
	pkgfunc.InitAksk()
	pkglog.InitLog(logimpl.NewLogImplSls(), "echo-svc")

	pkglog.Infov("EvtFngInitBegin",
		"AppVersion", version.GetAppVersion(),
		"GitCommit", version.GetShortGitCommit(),
		"BuildTime", version.GetBuildTimeStr())

	// 初始化 conf/dns
	pkgconf.InitConf(confimpl.LoadConfigOts)
	pkgddns.InitDdns(ddnsimpl.UpdateDns)

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
