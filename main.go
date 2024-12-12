package main

import (
	"context"
	"net"
	"net/http"
	"os"

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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	config, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	if config.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

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
		log.Fatal().Err(err).Msg("failed to create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migration up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(db db.Store, config utils.Config) {
	server, err := gapi.NewServer(db, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create gRPC server")
	}
	log.Info().Msg("gRPC server started")
	grpcLogger := grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return gapi.GrpcLogger(ctx, req, info, handler)
	})
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) // agar gprc client bisa tau ada apaaja di dalem grpc servernya

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Info().Msg("start gRPC server at " + listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}

func runGatewayServer(db db.Store, config utils.Config) {
	server, err := gapi.NewServer(db, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
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
		log.Fatal().Err(err).Msg("failed to register handler")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Info().Msg("start HTTP gateway server at " + listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func runHttpServer(db db.Store, config utils.Config) {
	server, err := api.NewServer(db, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
	log.Info().Msg("start HTTP server at " + config.HTTPServerAddress)
}
