package gapi

import (
	"fmt"

	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/pb"
	"github.com/ilhamgepe/simplebank/token"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/ilhamgepe/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          utils.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, config utils.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %v", err)
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
