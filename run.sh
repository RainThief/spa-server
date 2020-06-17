#!/bin/bash
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
set -e

# @todo gen ssl certs

docker build -f build/Dockerfile -t tmp .

docker run --init --rm -it -p 80:80 -p 443:443 -v $(pwd)/configs/config.default.yaml:/config.yml --name tmp tmp $@

# @todo update readme with where certs arety