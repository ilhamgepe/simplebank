package gapi

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/ilhamgepe/simplebank/db/sqlc"
	"github.com/ilhamgepe/simplebank/pb"
	"github.com/ilhamgepe/simplebank/utils"
	"github.com/ilhamgepe/simplebank/val"
	"github.com/ilhamgepe/simplebank/worker"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username: req.GetUsername(),
			Password: hashedPassword,
			FullName: req.GetFullName(),
			Email:    req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return s.taskDistributor.DistributeTaskVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %v", err)
			}
			return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
