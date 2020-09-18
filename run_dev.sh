#!/usr/bin/env bash
set -u

# Use this file to build and run the project for use in developing the application
# to force rebuild of go cache delete ./cache_built file

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

IMAGE="temp-spa-server"
CHECK_FILE="cache-built"
CACHE_IMAGE="$IMAGE-cache"

RETRY=${RETRY:-"false"}
sleep 2

if [ ! -f "$CHECK_FILE" ]; then

docker build -t "$CACHE_IMAGE" -f - . <<EOF
FROM golang:1.15 AS build
COPY . /app
WORKDIR /app
RUN go build -o spa-server /app/cmd/spa-server
EOF

    touch "$CHECK_FILE"
fi

docker build --build-arg baseImage="$CACHE_IMAGE" -f build/Dockerfile -t "$IMAGE" .
if [ $? -ne 0 ]; then
    if [ "$RETRY" == "false" ]; then
        rm "$CHECK_FILE"
        RETRY=true
        . ./run_dev.sh
    else
        exit 1
    fi
fi

docker stop temp-spa-server

docker run --init --rm -t --network=host -v $(pwd)/configs/config.default.yaml:/config.yml --name "$IMAGE" "$IMAGE" $@
