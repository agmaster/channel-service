#!/usr/bin/env bash

set -e

if [ ! -f getpkg.sh ]; then
    echo 'getpkg.sh must be run within its container folder' 1>&2
    exit 1
fi

OLDGOPATH="$GOPATH"
export GOPATH=`pwd`

if [ ! -d bin ]; then
    mkdir bin
fi
if [ ! -d pkg ]; then
    mkdir pkg
fi
if [ ! -d src ]; then
    mkdir src
fi


# go get -u -v github.com/lib/pq
# go get -u -v github.com/go-sql-driver/mysql
# go get -u -v github.com/gorilla/mux
# go get -u -v github.com/gorilla/sessions
# go get -u -v github.com/robfig/cron
# go get -u -v github.com/Sirupsen/logrus
go get -u -v github.com/julienschmidt/httprouter
go get -u -v gopkg.in/mgo.v2
go get -u -v gopkg.in/mgo.v2/bson
go get -u -v gopkg.in/olivere/elastic.v2

export GOPATH="$OLDGOPATH"
export PATH="$OLDPATH"

echo 'finished'