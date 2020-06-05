#!/bin/bash

#https://stackoverflow.com/questions/3760086/automatic-tagging-of-releases

git_tag () {

    VERSION=$(git describe --tags > /dev/null 2>&1)
    if [ $? == 128 ]; then
        VERSION="v0.0.0"
    fi

    #replace . with space so can split into an array
    VERSION_BITS=(${VERSION//./ })

    #get number parts and increase last one by 1
    VNUM1=${VERSION_BITS[0]}
    VNUM2=${VERSION_BITS[1]}
    VNUM3=${VERSION_BITS[2]}
    VNUM3=$((VNUM3+1))

    #create new tag
    echo "$VNUM1.$VNUM2.$VNUM3"

}

git_tag