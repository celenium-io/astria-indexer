-include .env
export $(shell sed 's/=.*//' .env)


indexer:
	cd cmd/indexer && go run . -c ../../configs/dipdup.yml

api:
	cd cmd/api && go run . -c ../../configs/dipdup.yml

generate:
	go generate -v ./internal/storage ./internal/storage/types ./pkg/node

lint:
	golangci-lint run

test:
	go test -p 8 -timeout 60s ./...

cover:
	go test ./... -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

api-docs:
	cd cmd/api && swag init --md markdown -parseDependency --parseInternal --parseDepth 1

ga:
	go generate -v ./internal/storage ./internal/storage/types ./pkg/node
	cd cmd/api && swag init --md markdown -parseDependency --parseInternal --parseDepth 1

license-header:
	update-license -path=./ -license=./HEADER

build:
	docker-compose up -d --build

.PHONY: indexer api generate test lint cover api-docs ga license-header build