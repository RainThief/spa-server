#!/bin/bash

# run this manually until gital ci works

set -uo pipefail

SCRIPTS_ROOT="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPTS_ROOT/include.sh"

IMAGE="registry.gitlab.com/martinfleming/spa-server"

docker pull $IMAGE:latest
NEW_TAG="$(git_tag)"
docker tag $IMAGE:latest $IMAGE:"${NEW_TAG/v/}"

git tag "$NEW_TAG"

git push --tags

docker push $IMAGE:"${NEW_TAG/v/}"
