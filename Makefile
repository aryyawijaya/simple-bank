DB_URL = postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable

pull-postgres:
	docker pull postgres:16.0-alpine3.18

start-postgres:
	docker \
		run \
		--name postgres16-alpine \
		--network bank-network \
		-p 5433:5432 \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=secretpassword \
		-d \
		postgres:16.0-alpine3.18

logs-postgres:
	docker logs -f postgres16-alpine

createdb:
	docker \
		exec \
		-it \
		postgres16-alpine \
		createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16-alpine dropdb simple_bank

migrateup-all:
	migrate \
		-path db/migration \
		-database "${DB_URL}" \
		-verbose \
		up

migrateup-1:
	migrate \
		-path db/migration \
		-database "${DB_URL}" \
		-verbose \
		up 1

migratedown-all:
	migrate \
		-path db/migration \
		-database "${DB_URL}" \
		-verbose \
		down

migratedown-1:
	migrate \
		-path db/migration \
		-database "${DB_URL}" \
		-verbose \
		down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./... -count=1

psql-console:
	docker exec -it postgres16-alpine psql -U root -d simple_bank

server:
	go run main.go

build-app-image:
	docker build -t simple-bank:latest .

start-app-dev:
	docker \
		run \
		--name simple-bank \
		-p 8080:8080 \
		-d \
		simple-bank:latest

start-app-prod:
	docker \
		run \
		--name simple-bank \
		--network bank-network \
		-p 8080:8080 \
		-e GIN_MODE=release \
		-e DB_SOURCE=postgresql://root:secretpassword@postgres16-alpine:5432/simple_bank?sslmode=disable \
		-d \
		simple-bank:latest

logs-app:
	docker logs -f simple-bank-api-1

compose-up-prod:
	docker compose -f compose.prod.yaml up -d

compose-down-prod:
	docker compose -f compose.prod.yaml down

compose-up-dev:
	docker compose -f compose.dev.yaml up -d

compose-down-dev:
	docker compose -f compose.dev.yaml down

delete-image:
	docker rmi simple-bank-api

logs-all:
	docker-compose logs -t -f

mock:
	mockgen \
		-package mockdb \
		-destination db/mock/store.go \
		github.com/aryyawijaya/simple-bank/db/sqlc Store

format:
	go fmt ./...

query-update:
	make sqlc mock

db-docs:
	dbdocs build doc/db.dbml

db-schema:
	dbml2sql doc/db.dbml --postgres -o doc/schema.sql

proto:
	rm -f pb/*.go
	protoc \
		--proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
    	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    	proto/*.proto

evans:
	evans --host localhost --port 8081 -r repl

.PHONY: \
	pull-postgres \
	start-postgres \
	logs-postgres \
	createdb \
	dropdb \
	migrateup-all \
	migrateup-1 \
	migratedown-all \
	migratedown-1 \
	sqlc \
	test \
	server \
	build-app-image \
	start-app \
	mock \
	format \
	query-update \
	db-docs \
	db-schema \
	proto \
	evans \
