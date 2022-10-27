package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"log"
	"net/http"

	"go-remix/config"
	"go-remix/handler"
)

func inject(d *dataSources, cfg config.Config) (*gin.Engine, error) {
	log.Println("注入数据源")

	router := gin.Default()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CorsOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}))

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
	router.Use(sessions.Sessions("go-remix", store))

	handler.InjectRouter(router, cfg)

	return router, nil
}
