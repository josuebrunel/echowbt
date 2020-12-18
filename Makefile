test:
	go test -v -cover ./...
lint:
	golint ./...
clean:
	go clean
debug:
	dlv test github.com/josuebrunel/echowbt
all: test lint
