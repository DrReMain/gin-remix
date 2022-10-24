package repository

import (
	"fmt"
	"go-remix/model"
	"net/smtp"
)

type mailRepository struct {
	username string
	password string
	origin   string
}

func NewMailRepository(username string, password string, origin string) model.MailRepository {
	return &mailRepository{
		username: username,
		password: password,
		origin:   origin,
	}
}

func (m *mailRepository) SendResetMail(email, token string) error {
	msg := "From: " + m.username + "\n" +
		"To: " + email + "\n" +
		"Subject: Reset Email\n\n" +
		fmt.Sprintf("<a href=\"%s/reset-password/%s\">Reset Password</a>", m.origin, token)

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", m.username, m.password, "smtp.gmail.com"), m.username, []string{email}, []byte(msg))

	return err
}
