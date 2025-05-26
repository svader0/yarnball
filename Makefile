.PHONY: build run clean test

build:
	go build -o bin/yarnball ./cmd/main.go

repl: build
	./bin/yarnball

clean:
	rm -rf bin

test:
	go test -v ./...