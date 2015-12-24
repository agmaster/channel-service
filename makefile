SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=channel-service

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`


.DEFAULT_GOAL: $(BINARY)
	

all: dependencies install
	
	
clean:
	go clean -i ./...
	go clean -i -race ./...
	
dependencies:
	go build -v ./...

$(BINARY): $(SOURCES)
	go build -o ${BINARY} main.go

generate: clean
	go install -v ./...
	go generate ./...


.PHONY: install
	
install: generate
	go install -v ./...


# .PHONY: clean
# clean:
# 	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
