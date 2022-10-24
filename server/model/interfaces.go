package model

import "context"

type RedisRepository interface {
	SetResetToken(ctx context.Context, id string) (string, error)
	GetIdFromToken(ctx context.Context, token string) (string, error)
}

type MailRepository interface {
	SendResetMail(email string, html string) error
}
