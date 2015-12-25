#!/bin/sh

# install Godep
go get -v github.com/tools/godep

godep version

echo "go get -v dependency packages"
go get -v github.com/astaxie/beego/logs
go get -v github.com/cactus/go-statsd-client/statsd
go get -v github.com/julienschmidt/httprouter
go get -v gopkg.in/mgo.v2
go get -v gopkg.in/mgo.v2/bson
go get -v gopkg.in/olivere/elastic.v2
