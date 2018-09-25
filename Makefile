GOPATH := $(GOPATH)

build:
	GOPATH=$(shell pwd); go build
run: build
	./obot $(TESTARGS)

