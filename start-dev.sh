#!/bin/sh

set -e

echo "update env vars"
source $HOME/.profile

echo "update dependencies"
go mod download

echo "run db migration"
migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"