package notification

import (
	"github.com/tgkzz/notification/pkg/logger"
	"github.com/tgkzz/notification/pkg/whatsapp"
	"log/slog"
)

type WhatsappService struct {
	log    *slog.Logger
	Client *whatsapp.WppService
}

func newWhatsappService(instanceId, authToken string, log *slog.Logger) (*WhatsappService, error) {
	wpp, err := whatsapp.New(instanceId, authToken)
	if err != nil {
		return nil, err
	}

	return &WhatsappService{Client: wpp, log: log}, nil
}

func (ws *WhatsappService) SendMessage(to, msg string) error {
	const op = "whatsapp.SendMessage"

	log := ws.log.With(
		slog.String("op", op),
		slog.String("to", to),
		slog.String("msg", msg),
	)

	log.Info("sending message")

	if err := ws.Client.ProcessMessage(to, msg); err != nil {
		log.Error("failed to send message", logger.Err(err))
		return err
	}

	return nil
}
