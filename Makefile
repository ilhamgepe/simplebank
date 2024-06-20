postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:16.3-alpine3.20
startpg:
	docker start postgres16
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb --username=root simple_bank
psql:
	docker exec -it postgres16 psql --username=root

sqlc:
	sqlc generate

reset:	down up
	clear
create:
	migrate create -ext sql -seq -dir db/migrations $(name)
up:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migrations -verbose up
down:
	migrate -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -path db/migrations -verbose down

test:
	go test -v -cover ./...