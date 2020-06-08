#!/bin/bash

set -euo pipefail

# in fetch mode it gets a commit not a branch
git checkout master

# ensure all tags are up to date (for fetch mode)
git fetch --all --tags

TAG=$(git tag | sort -V | tail -1)
if [ "$TAG" == "" ]; then
    TAG="v0.0.0"
fi

# strip preceeding "v" from tag
TAG="${TAG/v/}"

# get tag parts https://stackoverflow.com/questions/3760086/automatic-tagging-of-releases
TAG_BITS=(${TAG//./ })
VNUM1="${TAG_BITS[0]}"
VNUM2="${TAG_BITS[1]}"
VNUM3="${TAG_BITS[2]}"

# empty args do patch
if [ "$#" = "0" ]; then
    VNUM3=$((VNUM3+1))
fi

while [[ "$#" -gt 0 ]]
do
key="$1"

# bump version type based on arg passed
case $key in
    patch)
    VNUM3=$((VNUM3+1))
    shift
    ;;
    minor)
    VNUM2=$((VNUM2+1))
    VNUM3=0
    shift
    ;;
    major)
    VNUM1=$((VNUM1+1))
    VNUM2=0
    VNUM3=0
    shift
    ;;
    *)
    VNUM3=$((VNUM3+1))
    shift
    ;;
esac
done

NEW_TAG="$VNUM1.$VNUM2.$VNUM3"

GIT_URL="${CI_PROJECT_URL/https\:\/\//}"

git tag "v$NEW_TAG"
git push https://$CI_DEPLOY_USER:$ACCESS_TOKEN@$GIT_URL.git --tags

docker login $CI_REGISTRY -u $CI_REGISTRY_USER -p $CI_JOB_TOKEN
docker pull "$CI_REGISTRY_IMAGE":latest
docker tag "$CI_REGISTRY_IMAGE":latest "$CI_REGISTRY_IMAGE":"$NEW_TAG"
docker push "$CI_REGISTRY_IMAGE":"$NEW_TAG"
