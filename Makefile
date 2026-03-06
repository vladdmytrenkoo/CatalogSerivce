gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/product/v1/product_service.proto

docker-up:
	docker compose up -d

docker-down:
	docker compose down


test:
	go test ./...

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go