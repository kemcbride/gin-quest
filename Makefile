.PHONY: prod

prod:
	GIN_MODE=release go run main.go

dev:
	go run main.go
