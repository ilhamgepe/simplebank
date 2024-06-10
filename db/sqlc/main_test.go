package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries
var testDb *pgx.Conn

const (
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("closing database")
		if err := conn.Close(ctx); err != nil {
			log.Printf("failed to close database: %v\n", err)
		}
		log.Println("database closed")
	}()

	if err := conn.Ping(ctx); err != nil {
		log.Printf("failed to connect database: %v\n", err)
	}
	log.Println("database connected")
	testDb = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
