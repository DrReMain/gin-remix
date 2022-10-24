package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Domain         string `env:"DOMAIN"`
	Port           string `env:"PORT,required"`
	DatabaseUrl    string `env:"DATABASE_URL,required"`
	RedisUrl       string `env:"REDIS_URL,required"`
	SessionSecret  string `env:"SECRET,required"`
	SecretKey      string `env:"SECRET_KEY"`
	MailUser       string `env:"MAIL_USER"`
	MailPassword   string `env:"MAIL_PASSWORD"`
	CorsOrigin     string `env:"CORS_ORIGIN,required"`
	HandlerTimeOut int64  `env:"HANDLER_TIMEOUT,default=5"`
	MaxBodyBytes   int64  `env:"MAX_BODY_BYTES,default=4194304"` // default 4M
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}
	return
}
