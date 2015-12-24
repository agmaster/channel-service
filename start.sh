#!/bin/sh

export GOPATH=`pwd`

make clean
make
./channel-service & 


