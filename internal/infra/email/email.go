package emailInfra

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed templates/*.html
var templates embed.FS

type Template string

const (
	TemplateUserInvitation Template = "user_invitation"
)

const PyprHelloEmailId string = "Pypr <hello@piper.run.place>"

type EmailMessage struct {
	From    string
	To      string
	Subject string
	HTML    string
}

type EmailSender interface {
	Send(message EmailMessage) error
}

func (t Template) ReplaceTemplateVariables(variables map[string]string) (string, error) {
	file := fmt.Sprintf("templates/%s.html", t)

	content, err := templates.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("load email template %s: %w", t, err)
	}

	html := string(content)

	for key, value := range variables {
		html = strings.ReplaceAll(html, "{{"+key+"}}", value)
	}

	return html, nil
}
