package mails

import (
	mail "github.com/xhit/go-simple-mail/v2"
)

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

type C interface {
	Send(*Message) error
}

type Client struct{}

func (c *Client) Send(message *Message) error {
	msg := mail.NewMSG()
	msg.SetFrom(message.From)
	msg.AddTo(message.To)
	msg.SetSubject(message.Subject)
	msg.SetBody(mail.TextPlain, message.Body)
	err := msg.Send(smtpClient)

	return err
}
