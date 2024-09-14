package notification

import (
	"errors"
	"log/slog"
)

const (
	Whatsapp = "whatsapp"
	Telegram = "telegram"
)

type NotificationFactory interface {
	SendMessage(to, msg string) error
}

type NotifierService struct {
	WhatsappService *WhatsappService
}

func CreateNotifierService(instanceId, authToken string, log *slog.Logger) (NotifierService, error) {
	wppService, err := newWhatsappService(instanceId, authToken, log)
	if err != nil {
		return NotifierService{}, err
	}

	return NotifierService{
		WhatsappService: wppService,
	}, nil
}

func (n *NotifierService) GetNotificationService(name string) (NotificationFactory, error) {
	switch name {
	case Whatsapp:
		return n.WhatsappService, nil
	default:
		return nil, errors.New("unknown service name")
	}
}
