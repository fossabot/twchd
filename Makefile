BINARY_NAME	:=	bot
LDFLAGS		:=	-s -w

.PHONY: build test clean

build:
	go build -v -o $(BINARY_NAME) -ldflags "${LDFLAGS}"

test:
	go test -v ./...

clean:
	@go clean
	@rm -f $(BINARY_NAME)