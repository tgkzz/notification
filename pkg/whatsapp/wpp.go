package whatsapp

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

const UltraMsgBaseUrl = "https://api.ultramsg.com/"

type WppService struct {
	ultraMsgClient *resty.Client
}

func New(instanceId, authToken string) (*WppService, error) {
	if instanceId == "" || authToken == "" {
		return nil, errors.New("instance id and auth token should be given")
	}

	cli := resty.New().
		SetBaseURL(UltraMsgBaseUrl+instanceId).
		SetQueryParam("token", authToken)

	return &WppService{ultraMsgClient: cli}, nil
}

func (s *WppService) ProcessMessage(to, msg string) error {
	route := "/messages/chat"

	req, err := s.ultraMsgClient.R().
		SetQueryParams(map[string]string{
			"to":  to,
			"msg": url.QueryEscape(msg),
		}).Post(route)
	if err != nil {
		return err
	}

	if req.IsError() {
		return fmt.Errorf("could not send wpp message %s", req.Body())
	}

	return nil
}
