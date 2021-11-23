#!/usr/bin/env bash
set -eu

CI=${CI:-"false"}

IMAGE_NAME="ghcr.io/rainthief/spa-server/spa-server"

# Assume this script is in the src directory and work from that location
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd)/../"


pushd "$PROJECT_ROOT"
docker build --pull -t "$IMAGE_NAME" -f build/Dockerfile .

if [ "$CI" == "true" ]; then

    docker login "$DOCKER_REG" -u "$DOCKER_USER" -p "$DOCKER_PASS"

    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    BRANCH_TAG="${GIT_BRANCH//\//--}"

    docker tag "$IMAGE_NAME" "$IMAGE_NAME:$BRANCH_TAG"

    docker push "$IMAGE_NAME:$BRANCH_TAG"

    if [ "$GIT_BRANCH" == "master" ]; then
        # shellcheck disable=SC1091
        source "./scripts/semver.sh"
        RELEASE_TAG="$(get_tag "patch")"
        docker tag "$IMAGE_NAME" "$IMAGE_NAME:$RELEASE_TAG"
        git tag "v$RELEASE_TAG"
        git push --tags
        docker push "$IMAGE_NAME"
    fi
fi
