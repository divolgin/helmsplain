.PHONY: build
build:
	go build -o bin/helmsplain helmsplain.go

test:
	go test -v ./...