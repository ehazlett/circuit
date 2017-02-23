CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
COMMIT=`git rev-parse --short HEAD`
APP=circuit
REPO?=ehazlett/$(APP)
TAG?=latest
DEPS=$(shell go list ./... | grep -v /vendor/)

all: build

build:
	@cd cmd/$(APP) && go build -v -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

build-static:
	@cd cmd/$(APP) && go build -v -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

release: image
	@docker push $(REPO):$(TAG)

install: build
	@sudo cp -f cmd/$(APP)/$(APP) /usr/local/bin/$(APP)
test: build
	@go test -v $(DEPS)

clean:
	@rm -rf cmd/$(APP)/$(APP)

.PHONY: build build-static release test clean
