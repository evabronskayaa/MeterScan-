package mail

import (
	"context"
	"github.com/wneessen/go-mail"
	"html/template"
	"io"
	"log"
)

var t *template.Template

type Client struct {
	client *mail.Client

	login string
}

func NewMailClient(server string, port int, login, password string) (*Client, error) {
	c, err := mail.NewClient(server, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(login), mail.WithPassword(password), mail.WithSSL())
	if err != nil {
		return nil, err
	}
	return &Client{
		client: c,
		login:  login,
	}, nil
}

func (client *Client) SendPlainMessage(subject, message string, to ...string) error {
	m := mail.NewMsg()
	if err := m.From(client.login); err != nil {
		return err
	}
	if err := m.To(to...); err != nil {
		return err
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, message)

	if err := client.client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (client *Client) SendHtmlMessage(subject, file string, data any, to ...string) error {
	m := mail.NewMsg()
	if err := m.From(client.login); err != nil {
		return err
	}
	if err := m.To(to...); err != nil {
		return err
	}

	m.Subject(subject)
	m.SetBodyWriter(mail.TypeTextHTML, func(wr io.Writer) (int64, error) {
		return 0, t.ExecuteTemplate(wr, file, data)
	})

	if err := client.client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (client *Client) Start() error {
	var err error
	t, err = template.ParseGlob("templates/*.gohtml")
	return err
}

func (client *Client) Shutdown(_ context.Context) {
	if err := client.client.Close(); err != nil {
		log.Printf("Mail client close: %v", err)
	}
}
