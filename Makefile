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

migrateup:
	migrate \
		-path db/migration \
		-database "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable" \
		-verbose \
		up

migratedown:
	migrate \
		-path db/migration \
		-database "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable" \
		-verbose \
		down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

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
	docker logs -f simple-bank

compose-up:
	docker compose up

logs-all:
	docker-compose logs -t -f

.PHONY: 
	pull-postgres \
	start-postgres \
	logs-postgres \
	createdb \
	dropdb \
	migrateup \
	migratedown \
	sqlc \
	test \
	server \
	build-app-image \
	start-app
