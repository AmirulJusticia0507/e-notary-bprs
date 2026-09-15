.PHONY: fmt vet test build run migrate-up migrate-down tidy

DB_URL ?= postgres://postgres:postgres@localhost:5432/bprs_enotary?sslmode=disable

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

run:
	go run ./cmd/api

tidy:
	go mod tidy

migrate-up:
	migrate -path ../db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ../db/migrations -database "$(DB_URL)" down 1
