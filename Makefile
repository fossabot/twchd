PROJECTNAME := 	$(shell basename "$(PWD)")d
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOFILES 	:= 	$(wildcard *.go)
LDFLAGS 	= 	-ldflags "-s -w"

build:
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

test:
	go test -v ./...

clean:
	@echo "Cleaning build cache and binaries..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME)

deps:
	@echo "Installing missing dependencies..."
	@go get

image: build
	docker build -t aded/$(PROJECTNAME):$(shell git describe --tags) .
