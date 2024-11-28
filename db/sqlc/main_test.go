package db

import (
	"context"
	"os"
	"testing"

	"github.com/ilhamgepe/simplebank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
// )

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../../")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	testDB.Config().MaxConnIdleTime = 5 * 60
	testDB.Config().MaxConnLifetime = 10 * 60
	testDB.Config().MaxConns = 30
	testDB.Config().MinConns = 2

	testQueries = New(testDB)
	// jalanin test dulu dan simpan hasilnya ke test
	test := m.Run()
	// close database connection setelah test selesai

	// exit dan status yang di dapatkan dari test
	os.Exit(test)
}
