package notification

import "github.com/tgkzz/notification/pkg/whatsapp"

type WhatsappService struct {
	Client *whatsapp.WppService
}

func newWhatsappService(instanceId, authToken string) (*WhatsappService, error) {
	wpp, err := whatsapp.New(instanceId, authToken)
	if err != nil {
		return nil, err
	}

	return &WhatsappService{Client: wpp}, nil
}

func (ws *WhatsappService) SendMessage(to, msg string) error {
	return ws.Client.ProcessMessage(to, msg)
}
