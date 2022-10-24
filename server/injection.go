package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"go-remix/config"
	"go-remix/handler"
	"go-remix/model"
	"go-remix/repository"
	"go-remix/service"
	"log"
	"net/http"
	"time"
)

func inject(d *dataSources, cfg config.Config) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)

	redisRepository := repository.NewRedisRepository(d.RedisClient)
	mailRepository := repository.NewMailRepository(cfg.MailUser, cfg.MailPassword, cfg.CorsOrigin)

	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		RedisRepository: redisRepository,
		MailRepository:  mailRepository,
	})

	router := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CorsOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	router.Use(c)

	redisURL := d.RedisClient.Options().Addr
	password := d.RedisClient.Options().Password
	store, err := redis.NewStore(10, "tcp", redisURL, password, []byte(cfg.SessionSecret))
	if err != nil {
		return nil, fmt.Errorf("could not initialize redis session store: %w", err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		Domain:   cfg.Domain,
		MaxAge:   60 * 60 * 24 * 7, // 7days
		Secure:   gin.Mode() == gin.ReleaseMode,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions(model.CookieName, store))

	// TODO: ws

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		TimeoutDuration: time.Duration(cfg.HandlerTimeOut) * time.Second,
		MaxBodyBytes:    cfg.MaxBodyBytes,
	})

	return router, nil
}
