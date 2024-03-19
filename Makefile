.PHONY: build-linux-aarch64 build-linux-x86_64

BINARY_NAME=nodebox-dashboard

build-linux-aarch64:
	env GOOS=linux GOARCH=arm64 go build -o bin/nodebox-dashboard-linux-aarch64 main.go

build-linux-x86_64:
	env GOOS=linux GOARCH=amd64 go build -o bin/nodebox-dashboard-linux-x86_64 main.go

all: build-linux-aarch64 build-linux-x86_64

clean:
	go clean
	rm -rf bin/*

dev:
	mkdir -p /tmp/${BINARY_NAME}
	air .
