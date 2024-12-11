.PHONY: run build test swagger docker-build docker-run docker-migrate docker-clean docker-all

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

docker-build:
	docker compose build

docker-run:
	docker compose up app

docker-migrate:
	docker compose up migration

docker-clean:
	docker compose down -v

# Run everything in docker
docker-all: docker-build
	docker compose up -d elasticsearch
	sleep 10
	docker compose up migration
	docker compose up -d app
