PROJECT=github.com/fulldump/golax

GOCMD=go

.PHONY: test clean dependencies setup example coverage

all:	test

clean:
	rm -fr src

setup:
	mkdir -p src/$(PROJECT)
	rmdir src/$(PROJECT)
	ln -s ../../.. src/$(PROJECT)

test:
	$(GOCMD) version
	$(GOCMD) env
	$(GOCMD) test $(PROJECT)

example:
	$(GOCMD) install $(PROJECT)/example

coverage:
	$(GOCMD) test ./src/github.com/fulldump/goconfig -cover -covermode=count -coverprofile=coverage.out; \
	$(GOCMD) tool cover -html=coverage.out
