PROJECTNAME := 	$(shell basename "$(PWD)")d
GOBASE 		:= 	$(shell pwd)
GOBIN 		:= 	$(GOBASE)/bin
GOFILES 	:= 	$(filter-out assets_generate.go,$(wildcard *.go))
LDFLAGS 	= 	-ldflags "-s -w"

build: generate
	@echo "Building binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

test:
	go test -v ./...

clean:
	@echo "Cleaning build cache and binaries..."
	@go clean
	@rm -f $(GOBIN)/$(PROJECTNAME)
	@rm -f assets.go

image: build
	docker build -t aded/$(PROJECTNAME):$(shell git describe --tags) .

generate: assets_generate.go
	@echo "Embedding statics..."
	@go run $<

zsh-completion: $(GOBASE)/tools/_twchd
	@echo "ZSH Completion installing..."
	@cp $(GOBASE)/tools/_twchd /usr/share/zsh/site-functions
