PROJECTNAME := 	$(shell basename "$(PWD)")
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOFILES 	:= 	$(wildcard *.go)
LDFLAGS 	= 	-ldflags "-s -w"

build:
	@echo "Building binary..."
	go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

test:
	go test -v ./...

clean:
	@echo "Cleaning build cache and binaries..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME)

deps:
	@echo "Installing missing dependencies..."
	@go get
