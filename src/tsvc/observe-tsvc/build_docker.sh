#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o z-docker-dir/observe-tsvc.bin ./*.go

cd z-docker-dir
docker rmi registry.cn-hangzhou.aliyuncs.com/k8sns/observe-tsvc:v0.4
docker build -t registry.cn-hangzhou.aliyuncs.com/k8sns/observe-tsvc:v0.4 .
#docker push registry.cn-hangzhou.aliyuncs.com/k8sns/observe-tsvc:v0.4
