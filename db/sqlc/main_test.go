package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Config().MaxConnIdleTime = 5 * 60
	conn.Config().MaxConnLifetime = 10 * 60
	conn.Config().MaxConns = 30
	conn.Config().MinConns = 2

	testQueries = New(conn)
	// jalanin test dulu dan simpan hasilnya ke test
	test := m.Run()
	// close database connection setelah test selesai

	// exit dan status yang di dapatkan dari test
	os.Exit(test)
}
