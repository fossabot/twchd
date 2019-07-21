PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOSOURCE 	:= 	$(filter-out $(wildcard *_generate.go),$(wildcard *.go))
VERSION		:=	$(shell git describe --tags)
LDFLAGS 	:= 	-ldflags "-s -w"
CGO_ENABLED	:=	0
USER		:=	aded

export CGO_ENABLED

.PHONY: build clean image image-db build-arm generate env

build: generate
	@echo "Building binary..."
	go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

clean:
	@echo "Cleaning build cache and binaries..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME) $(GOBIN)/$(PROJECTNAME)-armv7

image: build env
	@echo "Building docker images..."
	docker build -t ${USER}/$(PROJECTNAME):$(VERSION) .

image-db: env
	docker build -t ${USER}/postgres-twchd:$(VERSION) -f Dockerfile-db .

build-arm:
	@echo "Building binary for ARMv7..."
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-armv7 $(GOSOURCE)

generate: assets
	@echo "Embedding statics..."
	@go run assets_generate.go

env:
	echo "export TAG=${VERSION}" > .env
