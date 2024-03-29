
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=channel-service

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`


.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} main.go

.PHONY: clean
clean:
	go clean -i .