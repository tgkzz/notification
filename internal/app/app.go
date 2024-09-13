package app

import (
	"context"
	"github.com/tgkzz/notification/internal/app/kafka"
	"log/slog"
)

type App struct {
	KafkaApp *kafka.App
}

func New(ctx context.Context, log *slog.Logger) (*App, error) {
	k, err := kafka.New(log)
	if err != nil {
		return nil, err
	}

	return &App{KafkaApp: k}, nil
}
