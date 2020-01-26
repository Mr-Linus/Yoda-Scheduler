all: local

local:
	GOOS=linux GOARCH=amd64 go build  -o=yoda-scheduler ./cmd/scheduler

build:
	sudo docker build --no-cache . -t registry.cn-hangzhou.aliyuncs.com/geekcloud/yoda-scheduler

push:
	sudo docker push registry.cn-hangzhou.aliyuncs.com/geekcloud/yoda-scheduler

format:
	sudo gofmt -l -w .
clean:
	sudo rm -f yoda-scheduler
