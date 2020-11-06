#!/usr/bin/env bash
set -u

# Use this file to build and run the project for use in developing the application
# to force rebuild of go cache delete ./cache_built file

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# @todo tidy this shit
HEALTH_CHECK_PORT="$(sed -nr 's/(.*)healthCheckPort:( )*(.*)/\3/p' configs/config.default.yaml)"
if [ "$HEALTH_CHECK_PORT" == "" ]; then
    HEALTH_CHECK_PORT="$(sed -nr 's/(.*)healthCheckDefaultPort =( )*(.*)/\3/p' internal/server/server.go)"
fi
echo $HEALTH_CHECK_PORT
CURRENT_HEALTH_CHECK_PORT="$(sed -nr 's/(.*)localhost:([0-9]+)(.*)/\2/p' build/Dockerfile)"
echo $CURRENT_HEALTH_CHECK_PORT
if [ "$CURRENT_HEALTH_CHECK_PORT" != "$HEALTH_CHECK_PORT" ]; then
    echo $HEALTH_CHECK_PORT
    echo "here"
    cp -f build/Dockerfile build/Dockerfile.bak
    sed -r "s/localhost:[0-9]+/localhost:$HEALTH_CHECK_PORT/g" build/Dockerfile > build/Dockerfile.new
    mv build/Dockerfile.new build/Dockerfile
fi
# exit
# # echo $DEFAULT_HEALTH_CHECK_PORT

# exit

# # @todo pushd

# # HEALTH_CHECK_PORT=$(echo $HEALTH_CHECK_PORT | xargs)
# cp -f build/Dockerfile build/Dockerfile.bak
# sed -r "s/http:\/\/localhost:[0-9]+/localhost:$HEALTH_CHECK_PORT/g" build/Dockerfile > build/Dockerfile.new
# mv build/Dockerfile.new build/Dockerfile


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
