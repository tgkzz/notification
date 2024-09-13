package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tgkzz/notification/internal/service/notification"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const topic = "notifications"

type Consumer struct {
	notification notification.NotifierService
	consumer     *kafka.Consumer
	log          *slog.Logger
}

func NewConsumer(instanceId, authToken, kafkaServer string, log *slog.Logger) (*Consumer, error) {
	n, err := notification.CreateNotifierService(instanceId, authToken)
	if err != nil {
		return nil, err
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"group.id":          "kafka-go-getting-started",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err = c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		notification: n,
		consumer:     c,
		log:          log,
	}, nil
}

func (c *Consumer) Consume() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			// log it as well
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log it
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	if err := c.consumer.Close(); err != nil {
		return
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
