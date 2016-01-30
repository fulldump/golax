PROJECT=golax
GOPATH=$(shell pwd)/_vendor
GOBIN=$(GOPATH)/bin
GOPKG=$(GOPATH)/pkg
GO=go
GOCMD=GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GO)

.DEFAULT_GOAL := build_example

.PHONY: all build clean dependencies setup

all: build

clean:
	rm -fr _vendor

setup:
	mkdir -p _vendor/src/github.com/fulldump
	ln -s ../../../.. _vendor/src/github.com/fulldump/golax
	ln -s ../../example _vendor/src/example

test:
	$(GOCMD) test ./...

dependencies:
	$(GOCMD) get $(PROJECT)

build_example: clean setup
	$(GOCMD) install example
