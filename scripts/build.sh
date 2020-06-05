#!/bin/bash

set -uo pipefail

SCRIPTS_ROOT="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPTS_ROOT/git_tag.sh"

TAG="$CI_REGISTRY_IMAGE":"${CI_COMMIT_REF_NAME//\//--}"
docker build -t $TAG .

if [ "$CI_COMMIT_REF_NAME" = "master" ]; then
    NEW_TAG="$(git_tag)"
    git tag $NEW_TAG
    git push --tags
    docker tag "$TAG" "$CI_REGISTRY_IMAGE":latest
    docker tag "$TAG" "$CI_REGISTRY_IMAGE":"$NEW_TAG"
fi

docker push "$CI_REGISTRY_IMAGE"