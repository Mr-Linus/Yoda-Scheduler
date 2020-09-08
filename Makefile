all: local

local:
	GOOS=linux GOARCH=amd64 go build  -o=yoda-scheduler ./cmd/scheduler

build:
	docker build --no-cache . -t registry.cn-hangzhou.aliyuncs.com/njupt-isl/yoda-scheduler:2.4

push:
	docker push registry.cn-hangzhou.aliyuncs.com/njupt-isl/yoda-scheduler:2.4

format:
	sudo gofmt -l -w .
clean:
	sudo rm -f yoda-scheduler
