GOPATH := $(GOPATH)

fmt:
	go fmt
build: fmt
	GOPATH=$(shell pwd); go build
run: build
	./obot
run-test: build
	./obot -p 10 < .test/test_file

