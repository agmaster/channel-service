#!/usr/bin/env bash

set -e

# if [ ! -f install.sh ]; then
#     echo 'install.sh must be run within its container folder' 1>&2
#     exit 1
# fi

CURDIR=`pwd`
# OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR:$CURDIR/../pkgs"

echo "$GOPATH"
#
# if [ ! -d bin ]; then
#     mkdir bin
# fi
# if [ ! -d pkg ]; then
#     mkdir pkg
# fi

#gofmt -w src

cd "$CURDIR/src/main" && go install -a

# export GOPATH="$OLDGOPATH"
# # OLDGOPATH="$GOPATH"
# export PATH="$OLDPATH"

echo 'finished'
