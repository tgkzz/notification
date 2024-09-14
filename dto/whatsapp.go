package dto

type WhatsappPayload struct {
	To        string `json:"to"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
}
