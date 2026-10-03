package emailInfra

import "fmt"

type Provider string

const (
	ProviderResend Provider = "RESEND"
)

type Factory struct {
	resend EmailSender
}

func NewFactory(
	resend EmailSender,
) *Factory {
	return &Factory{
		resend: resend,
	}
}

func (f *Factory) GetSender(provider Provider) (EmailSender, error) {
	switch provider {
	case ProviderResend:
		return f.resend, nil

	default:
		return nil, fmt.Errorf("unsupported email provider: %s", provider)
	}
}
