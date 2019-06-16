PROJECTNAME	:=	$(shell basename "$(PWD)")
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOTESTFILES	:= 	$(wildcard *_test.go)
GOSOURCE 	:= 	$(filter-out $(GOTESTFILES),$(wildcard *.go))
VERSION		:=	$(shell git describe --tags)
LDFLAGS 	:= 	-ldflags "-s -w"
DOCKERREPO	:=	127.0.0.1:1901

.PHONY: build clean image generate test dep

build: dep generate test
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOSOURCE)

clean:
	@echo "Cleaning build cache, binaries and assets..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME) $(GOBASE)/assets.go $(GOBIN)/$(PROJECTNAME)-armv7

image: build
	@echo "Building docker image..."
	docker build -t $(DOCKERREPO)/$(PROJECTNAME):$(VERSION) .

generate: assets/
	@echo "Embedding statics..."
	go-bindata -o assets.go assets/*

test:
	@echo "Running tests..."
	@go test ./...

coverage-report:
	go test -coverprofile /tmp/${PROJECTNAME}-general.cover
	grep -v 'assets.go' /tmp/${PROJECTNAME}-general.cover > /tmp/${PROJECTNAME}.cover
	go tool cover -html=/tmp/${PROJECTNAME}.cover

build-arm: dep generate test
	@echo "Building binary for ARMv7..."
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-armv7 $(GOSOURCE)

dep: Gopkg.toml Gopkg.lock
	@echo "Generating dependencies..."
	@dep ensure
