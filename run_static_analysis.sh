#!/bin/bash
set -e

# Assume this script is in the src directory and work from that location
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

docker run --rm -t \
    -v "$PROJECT_ROOT":/app \
    -w /app \
    golangci/golangci-lint:v1.27.0 \
    golangci-lint run -E golint -D gosimple
