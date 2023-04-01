package main

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	logFile *os.File
	logCh   chan int
)

func InitLog() {
	f, err := os.OpenFile("/tmp/observe-tsvc.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logFile = f
	writer := io.MultiWriter(f, os.Stdout)

	log.SetOutput(writer)
	log.Printf("InitLog ok, pid: %v", os.Getpid())

	logCh = make(chan int)

	// 定时打印，直到程序关闭
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ticker.C:
				log.Printf("log ticker")
			case <-logCh:
				log.Printf("stop log ticker")
				return
			}
		}
	}()
}

func DeInitLog() {
	close(logCh)

	log.Printf("DeInitLog ok")

	if logFile != nil {
		logFile.Close() // 可能会panic
		logFile = nil
	}
}
