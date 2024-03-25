BIN_DAEMON := "./bin/daemon"
BIN_CLIENT := "./bin/client"

DOCKER_IMG_DAEMON := "daemon:develop"
DOCKER_CONTAINER_DAEMON := "daemon"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	rm -rf internal/server/pb
	mkdir -p internal/server/pb

	protoc \
		--proto_path=api \
		--go_out=internal/server/pb \
		--go-grpc_out=internal/server/pb \
		api/*.proto

build-daemon:
	go build -v -o $(BIN_DAEMON) -ldflags "$(LDFLAGS)" ./cmd/daemon

build-client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

run-daemon: build-daemon
	$(BIN_DAEMON) -config ./config/config_daemon.toml

run-client: build-client
	$(BIN_CLIENT) -N 1 -M 5 -type 0

build-img-daemon:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_DAEMON) \
		-f build/Dockerfile .

run-img-daemon:
		docker run -d --name $(DOCKER_CONTAINER_DAEMON) \
    	-p 50051:50051 $(DOCKER_IMG_DAEMON)


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race -count=100 ./internal/...

.PHONY: generate lint test