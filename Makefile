BIN_DAEMON_LINUX := "./bin/linux/daemon"
BIN_CLIENT_LINUX := "./bin/linux/client"

BIN_DAEMON_FREEBSD := "./bin/freebsd/daemon"
BIN_CLIENT_FREEBSD := "./bin/freebsd/client"

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

build-daemon-linux:
	go build -v -o $(BIN_DAEMON_LINUX) -ldflags "$(LDFLAGS)" ./cmd/daemon

build-client-linux:
	go build -v -o $(BIN_CLIENT_LINUX) -ldflags "$(LDFLAGS)" ./cmd/client

build-daemon-freebsd:
	GOOS=freebsd go build -v -o $(BIN_DAEMON_FREEBSD) -ldflags "$(LDFLAGS)" ./cmd/daemon

build-client-freebsd:
	GOOS=freebsd go build -v -o $(BIN_CLIENT_FREEBSD) -ldflags "$(LDFLAGS)" ./cmd/client

run-daemon-linux: build-daemon-linux
	$(BIN_DAEMON_LINUX) -config ./config/config_daemon.toml

run-client-linux: build-client-linux
	$(BIN_CLIENT_LINUX) -N 1 -M 5 -type 0

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