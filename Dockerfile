# syntax = docker/dockerfile:experimental
FROM golang:1.12-alpine as build

RUN apk update && apk add make build-base git

ADD . /go/src/github.com/ehazlett/circuit
WORKDIR /go/src/github.com/ehazlett/circuit

ENV GO111MODULE=on
RUN --mount=type=cache,target=/root/.cache/go-build make build

FROM scratch as binary
COPY --from=build /go/src/github.com/ehazlett/circuit/bin/circuit /

FROM alpine:latest
COPY --from=build /go/src/github.com/ehazlett/circuit/bin/* /usr/local/bin/
