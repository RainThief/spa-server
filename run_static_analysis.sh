#!/bin/bash
set -e

# Assume this script is in the src directory and work from that location
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# if [[ ! -e .gomodules ]]; then
#     mkdir .gomodules
# fi


docker run --rm -t \
    -v "$PROJECT_ROOT":/app -u=$(id -u):$(id -g) \
    -w /app \
    -e XDG_CACHE_HOME=/tmp/.cache \
    golangci/golangci-lint:v1.27.0 \
    golangci-lint --color always run -E golint --skip-files test.go

docker run --rm -t -u=$(id -u):$(id -g) \
    -v "$PROJECT_ROOT":/app \
    -w /app \
    -e XDG_CACHE_HOME=/tmp/.cache \
    securego/gosec:v2.3.0 \
    -tests -severity LOW ./...

docker run --rm -t \
    -v "$PROJECT_ROOT"/go.sum:/app/go.sum \
    -w /app \
    sonatypecommunity/nancy:v0.3 \
    go.sum
