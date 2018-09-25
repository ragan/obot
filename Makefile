GOPATH := $(GOPATH)

build:
	GOPATH=$(shell pwd); go build
run: build
	./obot
run-test: build
	./obot -p 10 < .test/test_file

