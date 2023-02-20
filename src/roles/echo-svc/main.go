package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgconfig/configimpl"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgutil"
	"github.com/nysanier/fng/src/pkg/pkgvar"
	"github.com/nysanier/fng/src/pkg/version"
	"github.com/nysanier/fng/src/roles/echo-svc/router"
)

func main() {
	pkgvar.FnStartTime = getCRFC3339CstTimeStr()
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Printf("------- fng init begin -------")
	log.Printf("app version: %v.%v, build time: %v", version.AppVer, version.GetShortGitCommit(), version.GetBuildTimeStr())

	pkgenv.LoadEnv()
	configimpl.StartConfigUpdater()
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

	log.Printf("--- fng init end ---")
	var ch chan int
	<-ch
}

func getCRFC3339CstTimeStr() string {
	t := pkgfunc.GetCstNow()
	str := pkgfunc.GetRFC3339TimeStr(t)
	return str
}
