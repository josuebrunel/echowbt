test:
	go test -v -cover ./...
lint:
	golint ./...
clean:
	go clean
all: test lint
