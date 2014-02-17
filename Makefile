export GOPATH=$(shell pwd)

install:
	@go install baidu-pic
run:
	pkill baidu-pic
	@nohup ./bin/baidu-pic -port :8999 &
test:
	@./bin/baidu-pic -port :8999
