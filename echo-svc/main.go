package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := InitRouter()
	server := &http.Server{
		Addr: ":17080",
		Handler: r,
		ReadTimeout: 30*time.Second,
		WriteTimeout: 30*time.Second,
		MaxHeaderBytes: 2*1024*1024,
	}
	server.ListenAndServe()
}

func InitRouter() *gin.Engine {
	//gin.SetMode()
	r := gin.New()
	//v2 := r.Group("/v2")
	r.GET("/", Index)
	return r
}

const (
	Ver = "v1"
)

func Index(ctx *gin.Context) {
	t := time.Now()
	str := t.Format(time.RFC3339)
	ctx.String(http.StatusOK, "[%v] hello, echo-svc(%v)", str, Ver)
}