package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tgkzz/notification/internal/handler"
	"github.com/tgkzz/notification/internal/service/notification"
	"github.com/tgkzz/notification/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {

	a.Consume()

	return nil
}

// Consume TODO: maybe use transactions in order to exactly send message?
func (a *App) Consume() {
	const op = "kafka.Consume"

	l := a.log.With(slog.String("op", op))

	l.Info("starting to consume messages from server", slog.Any("kafka server", "localhost:52160"))

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			l.Warn("Caught signal of terminating", slog.Any("sig", sig))
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

	a.log.With(slog.String("op", op)).Info("stopping kafka server")

	if err := a.consumer.Close(); err != nil {
		panic(err)
	}
}
