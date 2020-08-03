# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=alertmanager-discord
GO_MAIN=cmd/alertmanager-discord/main.go

.PHONY: all vendor

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(GO_MAIN)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GORUN) $(GO_MAIN)
deps:
	$(GOMOD) download all

vendor:
	$(GOMOD) vendor

fmt:
	golangci-lint run --fix

upgrade:
	$(GOGET) -t -u ./...
