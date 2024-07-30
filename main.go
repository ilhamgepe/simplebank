package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ilhamgepe/simpleBank/api"
	db "github.com/ilhamgepe/simpleBank/db/sqlc"
	"github.com/ilhamgepe/simpleBank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Config().MaxConns = 4
	conn.Config().MinConns = 1
	conn.Config().MaxConnLifetime = time.Hour
	conn.Config().MaxConnIdleTime = 5 * time.Minute

	if err := conn.Ping(ctx); err != nil {
		log.Printf("failed to connect database: %v\n", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(fmt.Sprintf(":%s", config.ServerAddress)); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
