PROJECT=golax
GOPATH=$(shell pwd)/_vendor
GOBIN=$(GOPATH)/bin
GOPKG=$(GOPATH)/pkg
GO=go
GOCMD=GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GO)

.DEFAULT_GOAL := test

.PHONY: test clean dependencies setup example

clean:
	rm -fr _vendor

setup:
	mkdir -p _vendor/src/github.com/fulldump
	rm -fr _vendor/src/github.com/fulldump/golax
	ln -s ../../../.. _vendor/src/github.com/fulldump/golax
	ln -s ../../example _vendor/src/example

test: clean setup dependencies
	$(GOCMD) test ./...

dependencies:
	$(GOCMD) get example

example: clean setup dependencies
	$(GOCMD) install example
