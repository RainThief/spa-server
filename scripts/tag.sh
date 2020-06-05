#!/bin/bash

set -uo pipefail

SCRIPTS_ROOT="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPTS_ROOT/include.sh"

git checkout master

git tag "$(git_tag)"

git push http://martinfleming:-s8y6aM53wA1HiGEtYLb@gitlab.com/martinfleming/spa-server.git --tags
