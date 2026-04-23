.PHONY: prod dev air lint build

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4 fmt
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4 run --fix

build:
	go build ./...

prod:
	DOMAIN=kemcbride.noho.st GIN_MODE=release go run main.go

dev:
	go run main.go

air:
	mkdir .bin; go run github.com/air-verse/air@latest --build.cmd="go build -o .bin ./... " --build.entrypoint=".bin/gin-quest"
