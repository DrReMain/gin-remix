package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-remix/appo"
	"go-remix/middleware"
	"go-remix/model"
)

type Handler struct {
	userService  model.UserService
	MaxBodyBytes int64
}

type Config struct {
	R               *gin.Engine
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	_ = &Handler{
		MaxBodyBytes: c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, appo.NewNotFound("api", c.Request.RequestURI))
	})

	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(c.TimeoutDuration, appo.NewServiceUnavailable()))
	}

	ag := c.R.Group("api/account")
	ag.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, appo.NewSuccess("ok"))
	})
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("配置session错误: %v\n", err.Error())
	}
}
