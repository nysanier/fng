package main

import (
	"fmt"
	"log"
	rand "math/rand"
	"net/http"
	"os"
	"time"
)

var (
	server *http.Server
)

type Handler struct {
}

func (p *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	case "/liveness":
		writer.Write(nil) // ok
		DoMetricProbeCount(req.RequestURI)
		log.Printf("/liveness ok")
	case "/readiness":
		writer.Write(nil) // ok
		DoMetricProbeCount(req.RequestURI)
		log.Printf("/readiness ok")
	default:
		// do log and metric
		startTime := time.Now()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		hostName, _ := os.Hostname()
		str := fmt.Sprintf("[%v | %v] ok", hostName, os.Getpid())
		writer.Write([]byte(str)) // response 200
		DoMetricHttpRequest(startTime)
		log.Printf("biz-ServerHTTP ok, RequestURI=%v", req.RequestURI)
	}
}

func StartBiz() {
	rand.Seed(time.Now().UnixNano())

	go func() {
		server = &http.Server{Addr: ":17081", Handler: &Handler{}}
		server.ListenAndServe()
		if err := server.ListenAndServe(); err != nil {
			log.Printf("biz ListenAndServe fail, err=%v", err)
		}
	}()

	log.Printf("StartBiz ok")
}

func StopBiz() {
	server.Close()
	log.Printf("StopBiz ok")
}
