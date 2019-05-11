PROJECTNAME	:=	$(shell basename "$(PWD)")d
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOTESTFILES	:= 	$(wildcard *_test.go)
GOSOURCE 	:= 	$(filter-out $(GOTESTFILES),$(wildcard *.go))
VERSION		:=	$(shell git describe --tags)
LDFLAGS 	:= 	-ldflags "-s -w"

.PHONY: build clean image generate test

build: generate test
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

clean:
	@echo "Cleaning build cache, binaries and assets..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME) $(GOBASE)/assets.go

image: build
	@echo "Building docker image..."
	docker build -t $(USER)/$(PROJECTNAME):$(VERSION) .

generate: assets/
	@echo "Embedding statics..."
	go-bindata -o assets.go assets/*

test:
	@echo "Running tests..."
	@go test ./...