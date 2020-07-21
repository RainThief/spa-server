#!/usr/bin/env bash

# Use this file to build and run the project for use in developing the application
# to force rebuild of go cache delete ./cache_built file

# @todo use https://www.npmjs.com/package/npm-watch to auto reload

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
set -eu

IMAGE="temp-spa-server"
CHECK_FILE="cache-built"
CACHE_IMAGE="$IMAGE-cache"

if [ ! -f "$CHECK_FILE" ]; then

docker build -t "$CACHE_IMAGE" -f - . <<EOF
FROM golang:1.14 AS build
COPY . /app
WORKDIR /app
RUN go build -o spa-server /app/cmd/spa-server
EOF

    touch "$CHECK_FILE"

fi

docker build --build-arg baseImage="$CACHE_IMAGE" -f build/Dockerfile -t "$IMAGE" .

docker run --init --rm -it -p 80:80 -p 443:443 -v $(pwd)/configs/config.default.yaml:/config.yml --name "$IMAGE" "$IMAGE" $@
