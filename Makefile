.PHONY: docker-up docker-down migrate test run proto generate-mocks

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

migrate:
	@SPANNER_EMULATOR_HOST=localhost:9010 go run ./cmd/migrate

test:
	@go test ./...

run:
	@go run ./cmd/server

build:
	@go build -o bin/product-catalog-service ./cmd/server

proto:
	@protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product/v1/product_service.proto

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

generate-mocks:
	@mockery