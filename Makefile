BUILD_PATH=bin/msgbird
IMAGE=msgbird
CONTAINER=msgbird0

help:
	cat Makefile

init:
	brew install glide

init_linux:
	go get github.com/Masterminds/glide

deps:
	glide update && glide install

test:
	go test

build:
	go build -o ${BUILD_PATH}

build_linux:
	env GOOS=linux GOARCH=386 go build -o ${BUILD_PATH}

run:
	./${BUILD_PATH}

docker:
	docker build -t ${IMAGE} .
	docker images | grep ${IMAGE}

docker_run:
#   Selfnote for port expose: `-p host:container`
	docker run -itd -p 8080:8080 --name=${CONTAINER} --restart=always ${IMAGE} /bin/bash
	docker ps | grep ${CONTAINER}

docker_clean:
	docker stop ${CONTAINER}
	docker rm ${CONTAINER}
	docker rmi ${IMAGE} -f

start_macos: init deps test build run
start_linux: init_linux deps test build_linux run

example_call:
	curl -d "recipient=+3670...&originator=+3120...&message=Hello world!" -X POST http://127.0.0.1:8080/api/v1/message

cyclo:
	@go get -t github.com/fzipp/gocyclo
	@echo "List by cyclomatic complexities (over 2, exclude vendor dir):"
	@gocyclo -over 2 . | grep -v vendor
