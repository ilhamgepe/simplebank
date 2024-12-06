package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

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

	db := db.NewStore(pool)

	runGrpcServer(db, config)

}

func runGrpcServer(db db.Store, config utils.Config) {
	server, err := gapi.NewServer(db, config)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) //agar gprc client bisa tau ada apaaja di dalem grpc servernya

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}
func runHttpServer(db db.Store, config utils.Config) {
	server, err := api.NewServer(db, config)
	if err != nil {
		panic(err)
	}

	log.Fatal(server.Start(config.HTTPServerAddress))
}
