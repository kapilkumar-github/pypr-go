package emailInfra

import (
	"fmt"

	"github.com/resend/resend-go/v4"
)

type ResendSender struct {
	client *resend.Client
}

func NewResendSender(apiKey string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
	}
}

func (r *ResendSender) Send(message EmailMessage) error {
	params := &resend.SendEmailRequest{
		From:    message.From,
		To:      []string{message.To},
		Subject: message.Subject,
		Html:    message.HTML,
	}

	_, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("send email via resend: %w", err)
	}

	return nil
}
