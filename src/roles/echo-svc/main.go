package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nysanier/fng/src/pkg/pkgconfig"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgutil"
	"github.com/nysanier/fng/src/pkg/pkgvar"
	"github.com/nysanier/fng/src/pkg/version"
)

var (
	startTime = "program start time"
)

func main() {
	startTime = GetCstTimeStr()
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Printf("fng init begin")

	pkgenv.LoadEnv()
	pkgconfig.LoadConfig()

	// dev环境暂时不启动dnsUpdater
	if !pkgvar.IsDevEnv() {
		go pkgutil.RunDnsUpdater()
	}

	// Start Http Server
	r := InitRouter()
	server := &http.Server{
		Addr:           ":17080",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 2 * 1024 * 1024,
	}
	go server.ListenAndServe()

	log.Printf("fng init end")
	var ch chan int
	<-ch
}

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//v2 := r.Group("/v2")
	r.GET("/", Index)
	return r
}

const (
	Fn3339     = "Mon, 02 Jan 2006 3:04:05 PM"
	BodyFormat = `hello, echo-svc,
    current is:  %v
    remote addr: %v

> app version: %v.%v
> build time:  %v
> service ip:  %v
> start time:  %v`
)

func GetCstTimeStr() string {
	t := pkgfunc.GetCstNow()
	str := t.Format(Fn3339)
	return str
}

func Index(ctx *gin.Context) {
	cstTimeStr := GetCstTimeStr()
	remoteAddr := ctx.Request.RemoteAddr
	log.Printf("remote address: %v", remoteAddr)
	str := fmt.Sprintf(BodyFormat, cstTimeStr, remoteAddr,
		version.AppVer, version.GetShortGitCommit(), version.BuildTime, pkgutil.GetServiceIP(), startTime)
	ctx.String(http.StatusOK, str)
}
