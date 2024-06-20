package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDb *pgxpool.Pool

const (
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	pool.Config().MaxConns = 4
	pool.Config().MinConns = 1
	pool.Config().MaxConnLifetime = time.Hour
	pool.Config().MaxConnIdleTime = 5 * time.Minute

	if err := pool.Ping(ctx); err != nil {
		log.Printf("failed to connect database: %v\n", err)
	}
	log.Println("database connected")
	testDb = pool
	testQueries = New(pool)

	os.Exit(m.Run())
}
