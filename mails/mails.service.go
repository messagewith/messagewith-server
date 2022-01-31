package mails

import (
	"fmt"
	"messagewith-server/env"
)

var (
	client C
)

type service struct{}

func GetService(c C) *service {
	client = c
	return &service{}
}

func (s *service) SendResetPasswordToken(email string, token string) bool {
	err := client.Send(&Message{
		From:    fmt.Sprintf("Messagewith <%v>", env.SmtpEmail),
		To:      email,
		Subject: "Reset your password on Messagewith.app",
		Body:    fmt.Sprintf("Your reset password token is: %v", token),
	})

	if err != nil {
		return false
	}

	return true
}
