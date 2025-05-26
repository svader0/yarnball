build:
	go build -o bin/yarnball ./cmd/main.go

run: build
	./bin/yarnball