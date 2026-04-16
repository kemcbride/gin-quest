.PHONY: prod dev

prod:
	DOMAIN=kemcbride.noho.st GIN_MODE=release go run main.go

dev:
	go run github.com/air-verse/air@latest --build.cmd="go build -o .bin ./... " --build.entrypoint=".bin/gin-quest"
