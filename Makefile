VERSION ?= dev

.PHONY: build-linux-aarch64 build-linux-x86_64 clean

build-linux-aarch64:
	env GOOS=linux GOARCH=arm64 go build -o bin/nodebox-dashboard-$(VERSION)-linux-aarch64 main.go

build-linux-x86_64:
	env GOOS=linux GOARCH=amd64 go build -o bin/nodebox-dashboard-$(VERSION)-linux-x86_64 main.go

all: clean build-linux-aarch64 build-linux-x86_64

clean:
	rm -rf bin/
