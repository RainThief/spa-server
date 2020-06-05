#!/bin/bash

set -uo pipefail

SCRIPTS_ROOT="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPTS_ROOT/include.sh"

TAG="$CI_REGISTRY_IMAGE":"${CI_COMMIT_REF_NAME//\//--}"
docker build -t $TAG .

if [ "$CI_COMMIT_REF_NAME" = "master" ]; then
    NEW_TAG="$(git_tag)"
    docker tag "$TAG" "$CI_REGISTRY_IMAGE":latest
    docker tag "$TAG" "$CI_REGISTRY_IMAGE":"${NEW_TAG/v/}"
fi

docker login $CI_REGISTRY -u $CI_REGISTRY_USER -p $CI_JOB_TOKEN
docker push "$CI_REGISTRY_IMAGE"
