.PHONY: docker-up docker-down migrate test run proto

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	SPANNER_EMULATOR_HOST=localhost:9010 go run ./cmd/migrate

test:
	go test ./...

run:
	go run ./cmd/server

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product/v1/product_service.proto

tidy:
	go mod tidy