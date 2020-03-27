# syntax = docker/dockerfile:experimental
FROM golang:1.12-alpine as docker

RUN apk update && apk add make build-base git

ADD . /go/src/github.com/docker/caspian
WORKDIR /go/src/github.com/docker/caspian

ENV GO111MODULE=on
RUN --mount=type=cache,target=/root/.cache/go-build make build

FROM alpine:latest
COPY --from=docker /go/src/github.com/docker/caspian/bin/* /usr/local/bin/
