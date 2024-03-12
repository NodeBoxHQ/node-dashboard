
BINARY_NAME=nodebox-dashboard

build:
	go mod tidy
	go build -o bin/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf bin/*

dev:
	mkdir -p /tmp/${BINARY_NAME}
	air .