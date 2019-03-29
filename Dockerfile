FROM golang:1.12 AS builder

ARG VERSION=dev
ARG USER=heetch

WORKDIR /go/src/github.com/heetch/Darien-technical-test

COPY . .

RUN make build

FROM ubuntu:bionic

LABEL quay.expires-after=8w

RUN groupadd -r heetch && useradd --no-log-init -r -g heetch heetch
USER heetch

COPY --from=builder --chown=heetch:heetch /go/src/github.com/heetch/Darien-technical-test/bin/swiss-army-knife /bin/swiss-army-knife
