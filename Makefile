DB_URL=postgresql://root:root@localhost:5432/simple_bank?sslmode=disable

pg:
	@docker compose up -d postgres

migrate:
	@migrate create -ext sql -dir db/migration -seq $(word 2, $(MAKECMDGOALS))

migrateup:
	@migrate -path db/migration -database $(DB_URL) -verbose up

migrateup1:
	@migrate -path db/migration -database $(DB_URL) -verbose up

migratedown:
	@migrate -path db/migration -database $(DB_URL) -verbose down

migratedown1:
	@migrate -path db/migration -database $(DB_URL) -verbose down 1

dbreset:
	@migrate -path db/migration -database $(DB_URL) -verbose down -all
	@migrate -path db/migration -database $(DB_URL) -verbose up

sqlc:
	@sqlc generate

test:
	@go test -v -cover ./...

mock:
	@mockgen -package mockdb -destination db/mock/store.go github.com/ilhamgepe/simplebank/db/sqlc Store

server:
	@go run main.go


.PHONY: pg migrateup migrateup1 migrate migratedown migratedown1 sqlc test server mock