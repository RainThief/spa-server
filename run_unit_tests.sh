#!/usr/bin/env bash

# Assume this script is in the src directory and work from that location
PROJECT_ROOT=$(cd "$(dirname $0)" && pwd)

docker run --rm  -t \
    -u $(id -u):$(id -g) \
    --mount type=bind,source="${PROJECT_ROOT}",target=/app \
    -w /app \
    -e XDG_CACHE_HOME=/tmp/.cache \
    golang:1.15 go test -race -cover ./... "$@"
