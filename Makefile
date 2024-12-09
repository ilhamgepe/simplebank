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

proto:
	@rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb \
	--grpc-gateway_opt=paths=source_relative \
    proto/*.proto


.PHONY: pg migrateup migrateup1 migrate migratedown migratedown1 sqlc test server mock proto