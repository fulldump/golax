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
	mkdir -p _vendor/src
	ln -s ../.. _vendor/src/golax
	ln -s ../../example _vendor/src/example

test:
	$(GOCMD) test ./...

dependencies:
	$(GOCMD) get $(PROJECT)

build_example: clean setup dependencies
	$(GOCMD) install example;

build: clean setup dependencies
	for GOOS in "windows" "linux" "darwin"; do \
		for GOARCH in "386" "amd64"; do \
			echo "Building $$GOOS-$$GOARCH..."; \
			echo "GOOS=$$GOOS GOARCH=$$GOARCH $(GOCMD) build -o $(GOBIN)/$(PROJECT).$$GOOS.$$GOARCH"; \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GOCMD) build -o $(GOBIN)/$(PROJECT).$$GOOS.$$GOARCH; \
		done \
	done
	ls $(GOBIN)
