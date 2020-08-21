#!/bin/bash
#############################################################################
#
# A collection of useful bash functions to be used in scripts. To use the
# functions include the line
#
#    source ./scripts/include.sh
#
# in your script
#
#############################################################################


_pushd(){
    command pushd "$@" > /dev/null
}

_popd(){
    command popd "$@" > /dev/null
}

build_tmp_image() {
_pushd "${PROJECT_ROOT}"
docker build -t "$1" . -f-<<EOF
FROM golang:1.14-buster
WORKDIR /usr/app
COPY go.mod /usr/app/go.mod
RUN go mod download
EOF
_popd
}

get_image_name() {
    echo "tmp-$(basename -- $1)-image"
}

normalise_path() {
    # convert cygwin path
    if [ $(echo "$1" | grep cygdrive) ]; then
        echo "$1" | sed -r -e 's/\/cygdrive\/([a-z])/\1:/g'
        return
    fi
    echo "$1"
}

echo_colour() {
    colour=$2
    no_colour='\033[0m'
    echo -e "${colour}$1${no_colour}"
}

echo_warning(){
    magenta='\033[0;33;1m'
    echo_colour "$1" "${magenta}"
}

echo_success(){
    green='\033[0;32;1m'
    echo_colour "$1" "${green}"
}

echo_danger(){
  red='\033[0;31;1m'
  echo_colour "$1" "${red}"
}


echo_info(){
  cyan='\033[0;36;1m'
  echo_colour "$1" "${cyan}"
}
