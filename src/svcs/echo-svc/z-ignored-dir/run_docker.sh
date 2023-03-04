#!/bin/bash

#docker run -d -p 17080:17080 --env "fn_env:dev" --env "fn_aes_key=18072812467xxxxx" nysanier/echo-svc:v1.0.5
docker run -d -p 17080:17080 --env "fn_env:dev" --env "fn_aes_key=18072812467xxxxx" \
  --mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly nysanier/echo-svc:v1.0.5
