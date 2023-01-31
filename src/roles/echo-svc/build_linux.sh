#!/bin/bash

source const.sh

build_time=`date +%s`
git_commit=`git rev-parse HEAD`
app_ver_key="github.com/nysanier/fng/src/pkg/version.AppVer"
git_commit_key="github.com/nysanier/fng/src/pkg/version.GitCommit"
build_time_key="github.com/nysanier/fng/src/pkg/version.BuildTime"
echo "app_ver: ${app_ver}"
echo "git_commit: ${git_commit}"
echo "build_time: ${build_time}"
ldflags="-X '${app_ver_key}=${app_ver}' -X '${git_commit_key}=${git_commit}' -X '${build_time_key}=${build_time}'"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${ldflags}" -v -o docker-dir/echo-svc.bin main.go
