GOPATH := $(GOPATH)

fmt:
	go fmt
build: fmt
	GOPATH=$(shell pwd); go build
run: build
	./obot
test-run: build
	./obot -p 10 < .testing/test_file

