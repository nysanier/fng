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
	"github.com/nysanier/fng/src/pkg/version"
)

var (
	startTime = "program start time"
)

func main() {
	startTime = getCRFC3339CstTimeStr()
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Printf("fng init begin")

	pkgenv.LoadEnv()
	pkgconfig.StartConfigUpdater()
	pkgutil.StartDnsUpdater()

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
	r.Use(Recover)
	r.GET("/", Index)
	return r
}

// 防止意外宕机
func Recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover ok, err=%v", err)
		}
	}()
}

const (
	Fn3339     = "Mon, 02 Jan 2006 3:04:05 PM"
	BodyFormat = `hello, echo-svc,
    current is:  %v
    remote addr: %v

> app version: %v.%v
> build time:  %v
> start time:  %v
> service ip:  %v`
)

func getFn1123CstTimeStr() string {
	t := pkgfunc.GetCstNow()
	str := t.Format(Fn3339)
	return str
}

func getCRFC3339CstTimeStr() string {
	t := pkgfunc.GetCstNow()
	str := pkgfunc.GetRFC3339TimeStr(t)
	return str
}

func Index(ctx *gin.Context) {
	curTimeStr := getFn1123CstTimeStr()
	remoteAddr := ctx.Request.RemoteAddr
	log.Printf("remote address: %v", remoteAddr)
	str := fmt.Sprintf(BodyFormat, curTimeStr, remoteAddr,
		version.AppVer, version.GetShortGitCommit(), version.GetBuildTimeStr(), startTime, pkgutil.GetCurrentServiceIP())
	ctx.String(http.StatusOK, str)
}
