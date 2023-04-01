package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init observable, not biz
	InitLog()
	InitMetric()
	InitTrace()

	// start biz
	StartBiz()

	// wait exit signal, then stopping biz
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	sig := <-ch
	log.Printf("sig: %v, pid: %v", sig, os.Getpid())

	// stop biz
	StopBiz()

	// wait biz stopped
	log.Printf("waiting main coprocess exit ...")
	time.Sleep(time.Second * 3)
	log.Printf("main coprocess exited")

	// deinit observable, until biz stopped
	InitTrace()
	DeInitMetric()
	DeInitLog()
}
