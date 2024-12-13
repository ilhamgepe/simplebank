package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type TaskDistributor interface {
	DistributeTaskVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{client}
}

func (distributor *RedisTaskDistributor) CloseClient() error {
	if err := distributor.client.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close redis client")
		return err
	}
	log.Info().Msg("Redis client closed successfully.")
	return nil
}
