package handler

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tgkzz/notification/dto"
	"github.com/tgkzz/notification/internal/service/notification"
	"log/slog"
)

type BrokerHandler struct {
	notificationService *notification.NotifierService
	logger              *slog.Logger
}

type Handler interface {
	Gateway(req *kafka.Message) error
}

func New(logger *slog.Logger, notificationService *notification.NotifierService) Handler {
	return &BrokerHandler{logger: logger, notificationService: notificationService}
}

func (bh *BrokerHandler) Gateway(req *kafka.Message) error {
	bh.logger.Info("Received message from Kafka", slog.String("topic", *req.TopicPartition.Topic))

	switch string(req.Key) {
	case "whatsapp":
		return bh.handleWhatsappMsg(req)
	default:
		bh.logger.Warn("Unknown message key", slog.String("key", string(req.Key)))
		return nil
	}
}

func (bh *BrokerHandler) handleWhatsappMsg(msg *kafka.Message) error {

	log := bh.logger.With(
		slog.String("notification service", string(msg.Key)),
		slog.String("payload", string(msg.Value)),
	)

	log.Info("Handled message")

	var req dto.WhatsappPayload
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Error("Failed to unmarshall kafka msg  to whatsapp payload", slog.String("error", err.Error()))
		return err
	}

	srv, err := bh.notificationService.GetNotificationService(string(msg.Key))
	if err != nil {
		log.Error("Failed to get notification service", slog.String("error", err.Error()))
		return err
	}

	if err = srv.SendMessage(req.To, req.Message); err != nil {
		log.Error("Failed to send message", slog.String("error", err.Error()))
		return err
	}

	return nil
}
