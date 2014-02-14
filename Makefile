export GOPATH=$(shell pwd)

install:
	@go install baidu-pic
run:
	@./bin/baidu-pic
test:
	@./bin/baidu-pic
