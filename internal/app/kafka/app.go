package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tgkzz/notification/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const topic = "notifications"

// App TODO: make handler to redirect messages over there
type App struct {
	consumer *kafka.Consumer
	log      *slog.Logger
}

func New(log *slog.Logger) (*App, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:52160",
		"group.id":          "kafka-go-getting-started",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err = c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &App{consumer: c, log: log}, nil
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "kafka.App"

	go func() {
		a.Consume()
	}()

	return nil
}

func (a *App) Consume() {
	const op = "kafka.Consume"

	l := a.log.With(slog.String("op", op))

	l.Info("starting to consume messages from server %s", "localhost:52160")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			l.Warn("Caught signal %v: terminating", sig)
			run = false
		default:
			ev, err := a.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				l.Error("Failed to read message", logger.Err(err))
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
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
