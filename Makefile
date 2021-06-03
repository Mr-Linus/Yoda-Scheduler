all: local

local:
	GOOS=linux GOARCH=amd64 go build  -o=yoda-scheduler ./cmd/scheduler

build: local
	docker build --no-cache . -t registry.cn-hangzhou.aliyuncs.com/geekcloud/yoda-scheduler:1.0

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/geekcloud/yoda-scheduler:1.0

format:
	gofmt -l -w .
clean:
	rm -f yoda-scheduler
