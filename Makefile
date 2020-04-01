# Copyright (c) Evan Hazlett
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation
# files (the "Software"), to deal in the Software without
# restriction, including without limitation the rights to use, copy,
# modify, merge, publish, distribute, sublicense, and/or sell copies
# of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
# DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
# OR OTHER DEALINGS IN THE SOFTWARE.

COMMIT=$(shell git rev-parse HEAD | head -c 8)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)
NAMESPACE?=ehazlett
IMAGE_NAMESPACE?=$(NAMESPACE)
APP=circuit
REPO?=$(NAMESPACE)/$(APP)
TAG?=dev
BUILD?=-dev
PACKAGES=$(shell go list ./... | grep -v /vendor/)
IMAGE?=$(REPO):latest

all: build-docker

bindir:
	@mkdir -p bin

build: bindir server

server: bindir
	@cd cmd/circuit && go build -v -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT) -X github.com/$(REPO)/version.Build=$(BUILD)" -mod=vendor -o ../../bin/circuit

build-docker:
	@docker build -t ${IMAGE} .

install:
	@install bin/* /usr/local/bin/

test:
	@go test -v ./...

protos:
	@protobuild --quiet ${PACKAGES}

clean:
	@rm -rf bin

PHONY: build contrib clean
