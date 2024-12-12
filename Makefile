.PHONY: run build test swagger

run:
	go run cmd/api/main.go

build:
	go build -o api cmd/api/main.go

test:
	go test -v ./...

swagger:
	swag init -g cmd/api/main.go -o api/swagger

migrate:
	go run cmd/migration/main.go
