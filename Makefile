test:
	go test -v ./...
lint:
	golint ./...
clean:
	go clean
all: test lint
