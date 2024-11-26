DB_URL=postgresql://root:root@localhost:5432/simple_bank?sslmode=disable

pg:
	@docker compose up -d postgres

migrateup:
	@migrate create -ext sql -dir db/migration -seq $(word 2, $(MAKECMDGOALS))

migrate:
	@migrate -path db/migration -database $(DB_URL) -verbose up

migratedown:
	@migrate -path db/migration -database $(DB_URL) -verbose down
migratedown1:
	@migrate -path db/migration -database $(DB_URL) -verbose down 1

sqlc:
	@sqlc generate

.PHONY: pg migrateup migrate migratedown sqlc