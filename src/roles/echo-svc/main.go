package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nysanier/fng/src/pkg/version"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	r := InitRouter()
	server := &http.Server{
		Addr:           ":17080",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 2 * 1024 * 1024,
	}
	server.ListenAndServe()
}

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//v2 := r.Group("/v2")
	r.GET("/", Index)
	return r
}

const (
	Lnatian3339 = "Mon, 02 Jan 2006 3:04:05 PM"
	BodyFormat  = `hello, echo-svc, current is
    %v

> app version: %v
> git commit: %v
> build time: %v`
)

func GetCstTimeStr() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("time.LoadLocation fail: err=%v", err)
		loc = time.UTC
	}

	t := time.Now().In(loc)
	str := t.Format(Lnatian3339)
	return str
}

func Index(ctx *gin.Context) {
	cstTimeStr := GetCstTimeStr()
	remoteAddr := ctx.Request.RemoteAddr
	log.Printf("remote address: %v", remoteAddr)
	ctx.String(http.StatusOK, BodyFormat, cstTimeStr, version.AppVer, version.GetShortGitCommit(), version.BuildTime)
}
