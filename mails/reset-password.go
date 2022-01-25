package mails

import (
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"messagewith-server/env"
)

func SendResetPasswordToken(email string, token string) bool {
	msg := mail.NewMSG()
	msg.SetFrom(fmt.Sprintf("Messagewith <%v>", env.SmtpEmail))
	msg.AddTo(email)
	msg.SetSubject("Reset your password on Messagewith.app")
	msg.SetBody(mail.TextPlain, fmt.Sprintf("Your reset password token is: %v", token))

	err := msg.Send(smtpClient)

	if err != nil {
		return false
	}

	return true
}
