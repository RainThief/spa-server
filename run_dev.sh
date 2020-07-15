#!/usr/bin/env bash

Use this file to build and run the project for use in developing the application
# @todo use https://www.npmjs.com/package/npm-watch to auto reload

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
set -eu

IMAGE="temp-spa-server"

docker build -f build/Dockerfile -t "$IMAGE" .

docker run --init --rm -it -p 80:80 -p 443:443 -v $(pwd)/configs/config.default.yaml:/config.yml --name "$IMAGE" "$IMAGE" $@
