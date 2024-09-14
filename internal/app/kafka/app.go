package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tgkzz/notification/internal/handler"
	"github.com/tgkzz/notification/internal/service/notification"
	"github.com/tgkzz/notification/pkg/logger"
	"log/slog"
)

const Topic = "notifications"

type App struct {
	consumer *kafka.Consumer
	log      *slog.Logger
	handler  handler.Handler
}

func New(log *slog.Logger, wppInstance, wppToken string) (*App, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:51807",
		"group.id":          "kafka-go-getting-started",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err = c.SubscribeTopics([]string{Topic}, nil); err != nil {
		return nil, err
	}

	n, err := notification.CreateNotifierService(wppInstance, wppToken, log)
	if err != nil {
		return nil, err
	}

	h := handler.New(log, &n)

	return &App{consumer: c, log: log, handler: h}, nil
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) run(ctx context.Context) error {

	a.Consume(ctx)

	return nil
}

// Consume TODO: maybe use transactions in order to exactly send message?
// Consume TODO: realize parallel processing of messages from kafka
func (a *App) Consume(ctx context.Context) {
	const op = "kafka.Consume"

	l := a.log.With(slog.String("op", op))

	l.Info("starting to consume messages from server", slog.Any("kafka server", "localhost:52160"))

	run := true
	for run {
		select {
		case <-ctx.Done():
			l.Warn("Context canceled, stopping consumer")
			run = false
		default:
			ev, err := a.consumer.ReadMessage(-1)
			if err != nil {
				l.Error("Failed to read message", logger.Err(err))
				continue
			}
			if err = a.handler.Gateway(ev); err != nil {
				l.Error("failed to read msg", logger.Err(err))
				return
			}
			if _, err = a.consumer.CommitMessage(ev); err != nil {
				l.Error("failed to commit msg", logger.Err(err))
			}
		}
	}
}

func (a *App) Stop() {
	const op = "kafka.Stop"

	a.log.With(slog.String("op", op)).Info("stopping kafka consumer")

	if err := a.consumer.Close(); err != nil {
		panic(err)
	}
}
