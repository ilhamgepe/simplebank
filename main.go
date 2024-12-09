package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ilhamgepe/simplebank/api"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/gapi"
	"github.com/ilhamgepe/simplebank/pb"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
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

	// run migration
	runDBMigration(config.MigrationURL, config.DBSource)

	db := db.NewStore(pool)

	go runGatewayServer(db, config)
	runGrpcServer(db, config)
}

func runDBMigration(migrationUrl, dbSource string) {
	m, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migration up: %v\n", err)
	}

	log.Println("db migrated successfully")
}

func runGrpcServer(db db.Store, config utils.Config) {
	server, err := gapi.NewServer(db, config)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) // agar gprc client bisa tau ada apaaja di dalem grpc servernya

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	fmt.Printf("start gRPC server at %s\n", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func runGatewayServer(db db.Store, config utils.Config) {
	server, err := gapi.NewServer(db, config)
	if err != nil {
		log.Fatalf("failed to create server: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("failed to register handler server: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	fmt.Printf("start HTTP gateway server at %s\n", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start HTTP gateway server: %v\n", err)
	}
}

func runHttpServer(db db.Store, config utils.Config) {
	server, err := api.NewServer(db, config)
	if err != nil {
		log.Fatalf("failed to create server: %v\n", err)
	}

	log.Fatal(server.Start(config.HTTPServerAddress))
}
