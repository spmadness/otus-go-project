FROM golang:1.21 as build

ENV BIN_FILE /opt/daemon/daemon-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/daemon/*

FROM ubuntu:20.04

ENV BIN_FILE "/opt/daemon/daemon-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/daemon/config.toml
COPY ./config/config_daemon.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}

