PROJECTNAME	:=	$(shell basename "$(PWD)")d
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOFILES 	:= 	$(filter-out assets_generate.go,$(wildcard *.go))
VERSION		:=	$(shell git describe --tags)
LDFLAGS 	:= 	-ldflags "-s -w"

.PHONY: build clean image generate

build: generate
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

clean:
	@echo "Cleaning build cache, binaries and assets..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME) $(GOBASE)/assets.go

image: build
	@echo "Building docker image..."
	docker build -t $(USER)/$(PROJECTNAME):$(VERSION) .

generate: $(GOBASE)/assets_generate.go 
	@echo "Embedding statics..."
	@go run $<
