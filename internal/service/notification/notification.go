package notification

import (
	"errors"
	"github.com/tgkzz/notification/internal/config"
)

const (
	Whatsapp = "whatsapp"
	Telegram = "telegram"
)

type NotificationFactory interface {
	SendMessage(to, msg string) error
}

type NotifierService struct {
	WhatsappConfig config.WhatsappConfig
}

func CreateNotifierService() (NotifierService, error) {
	return NotifierService{}, nil
}

func (n *NotifierService) GetNotificationService(name string) (NotificationFactory, error) {
	switch name {
	case Whatsapp:
		return n.createWhatsappService()
	default:
		return nil, errors.New("unknown service name")
	}
}

func (n *NotifierService) createWhatsappService() (NotificationFactory, error) {
	srv, err := NewWhatsappService(n.WhatsappConfig.InstanceId, n.WhatsappConfig.AuthToken)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
