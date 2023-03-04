#!/bin/bash

sh build_linux.sh

source const.sh

cd z-docker-dir
docker rmi nysanier/echo-svc:${app_ver}
docker build -t nysanier/echo-svc:${app_ver} .
#docker push nysanier/echo-svc:${app_ver} # 通过docker desktop来推送
