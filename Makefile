.PHONY: prod dev

prod:
	DOMAIN=kemcbride.noho.st GIN_MODE=release go run main.go

dev:
	go run main.go
