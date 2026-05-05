.PHONY: build run clean test

build:
	go build -o yarnball

repl: build
	./yarnball

clean:
	rm -rf bin

test:
	go test -v ./...