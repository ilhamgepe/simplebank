package main

import (
	"context"
	"log"

	"github.com/ilhamgepe/simplebank/api"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	pool.Config().MaxConnIdleTime = 5 * 60
	pool.Config().MaxConnLifetime = 10 * 60
	pool.Config().MaxConns = 30
	pool.Config().MinConns = 2

	db := db.NewStore(pool)

	server, err := api.NewServer(db, config)
	if err != nil {
		panic(err)
	}

	log.Fatal(server.Start(config.ServerAddress))
}
