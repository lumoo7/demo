package lib

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

type Email struct {
	SmtpHost string
	SmtpPort string
	Password string
	From     string
	To       string
	Cc       string
	Subject  string
	Text     []byte
	HTML     []byte
}

type IEmail interface {
	SendEmail(email *Email) error
	WriteText(msg string)
	InitEmail() *Email
}

func (e *Email) SendEmail(em *email.Email) error {
	return em.Send(
		strings.Join([]string{e.SmtpHost, e.SmtpPort}, ":"),
		smtp.PlainAuth("", e.From, e.Password, e.SmtpHost),
	)
}

func (e *Email) InitEmail() *email.Email {
	newEmail := email.NewEmail()
	newEmail.From = e.From
	newEmail.To = []string{e.To}
	newEmail.Subject = e.Subject
	newEmail.Text = e.Text
	return newEmail
}

func (e *Email) WriteText(msg string) {
	e.Text = []byte(msg)
}

func NewEmail(SmtpHost, SmtpPort, Password, From, To, Cc, Subject string, Text, HTML []byte) *Email {
	return &Email{
		SmtpHost: SmtpHost,
		SmtpPort: SmtpPort,
		Password: Password,
		From:     From,
		To:       To,
		Cc:       Cc,
		Subject:  Subject,
		Text:     Text,
		HTML:     HTML,
	}
}
