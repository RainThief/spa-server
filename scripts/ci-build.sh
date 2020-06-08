#!/bin/bash

set -uo pipefail

TAG="$CI_REGISTRY_IMAGE":"${CI_COMMIT_REF_NAME//\//--}"
docker build -t $TAG .

if [ "$CI_COMMIT_REF_NAME" = "master" ]; then
    docker tag "$TAG" "$CI_REGISTRY_IMAGE":latest
fi

docker login $CI_REGISTRY -u $CI_REGISTRY_USER -p $CI_JOB_TOKEN
docker push "$CI_REGISTRY_IMAGE"
