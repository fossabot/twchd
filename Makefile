PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOSOURCE 	:= 	$(wildcard *.go)
VERSION		:=	$(shell git describe --tags)
LDFLAGS 	:= 	-ldflags "-s -w"


.PHONY: build clean image build-arm

build:
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

clean:
	@echo "Cleaning build cache and binaries..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME) $(GOBIN)/$(PROJECTNAME)-armv7

image: build
	@echo "Building docker image..."
	docker build -t ${USER}/$(PROJECTNAME):$(VERSION) .

build-arm:
	@echo "Building binary for ARMv7..."
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-armv7 $(GOSOURCE)
