# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: wormholes all test clean

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

wormholes:
	go build -o $(GOBIN)/wormholes ./cmd/wormholes
	@echo "Done building."
	@echo "Run \"$(GOBIN)/wormholes\" to launch wormholes."

all:
	$(GORUN) build/ci.go install

test: all
	$(GORUN) build/ci.go test

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

