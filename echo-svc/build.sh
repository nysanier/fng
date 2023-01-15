#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o v1/echo-svc.bin main.go
