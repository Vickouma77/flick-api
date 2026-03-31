package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// templateFS holds embedded email templates from the templates directory.
//
//go:embed "templates"
var templateFS embed.FS

// Mailer wraps SMTP dialing configuration and sender identity.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// New initializes a Mailer with SMTP credentials and a default sender address.
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send renders an email template and sends the message to a recipient.
// The template must define "subject", "plainBody", and "htmlBody" blocks.
func (m *Mailer) Send(recipient, templateFile string, data any) error {
	// Parse the selected email template file from the embedded filesystem.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Render the subject block into a text buffer.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Render the plain-text body block for email clients without HTML support.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Render the HTML body block for richer email clients.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Build the outgoing message with headers and both body variants.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
