#!/bin/sh

set -e

echo "[INFO] Running migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "[INFO] Start server"
exec "$@"