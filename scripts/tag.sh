#!/bin/bash

set -euo pipefail

IMAGE="registry.gitlab.com/martinfleming/spa-server"

# run this manually until gital ci works
# Semver image
# Needs to pull latest
# Give option to minor major or patch
# Build and copy netrc and put in private repo semver project
# Env repo
# Copy in .netrc
# Make as own repo project then make second project to create image

git checkout master
git pull origin master

# delete current tags
git tag -l | xargs git tag -d

# fetch remote tags
git fetch --tags

# get git tag
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
echo $NEW_TAG
exit

docker pull "$IMAGE":latest
docker tag "$IMAGE":latest "$IMAGE":"$NEW_TAG"

git tag "v$NEW_TAG"

git push --tags

docker push "$IMAGE":"$NEW_TAG"
