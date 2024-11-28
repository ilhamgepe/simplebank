package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/ilhamgepe/simplebank/api"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	pool.Config().MaxConnIdleTime = 5 * 60
	pool.Config().MaxConnLifetime = 10 * 60
	pool.Config().MaxConns = 30
	pool.Config().MinConns = 2

	db := db.NewStore(pool)

	validate := validator.New(validator.WithRequiredStructEnabled())
	server := api.NewServer(db, validate)

	log.Fatal(server.Start(":8080"))
}
