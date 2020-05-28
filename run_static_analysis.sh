#!/bin/bash

# @todo fix this

# Assume this script is in the src directory and work from that location
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

export GO_DEV_IMAGE_TAG=registry.openffd.org/ffd/go-dev:1.0.74

set -e

# Copy the credentials required to checkout the go-utilities sub-projects
# required (e.g. logging, transport ...)
if [[ ! -f "${PROJECT_ROOT}/.netrc" ]];
then
  cp "${HOME}/.openffd.netrc" "${PROJECT_ROOT}/.netrc"
fi

mkdir -p .gomodules/pkg

docker run --rm -u $(id -u):$(id -g) -t \
  --mount type=bind,source="${PROJECT_ROOT}",target=/app \
  --mount type=bind,source="${PROJECT_ROOT}"/.netrc,target=/.netrc \
  -e GOPATH=/gopath -v "$(pwd)"/.gomodules:/gopath \
  -e PROJECT_ROOT=/app ${GO_DEV_IMAGE_TAG} /bin/go_lint_tool.sh

# If the credentials were copied from the users home directory delete them
# from the checkout
if [[ -f "${HOME}/.openffd.netrc" ]]; then
  rm -f "${PROJECT_ROOT}/.netrc"
fi
