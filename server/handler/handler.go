package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"go-remix/appo"
	"go-remix/middleware"
	"go-remix/model"
)

type Handler struct {
	MaxBodyBytes int64
	userService  model.UserService
}

type Config struct {
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
	R               *gin.Engine
	UserService     model.UserService
}

func NewHandler(c *Config) {
	h := &Handler{
		MaxBodyBytes: c.MaxBodyBytes,
		userService:  c.UserService,
	}

	// 不存在路由处理
	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, appo.NewNotFound("api", c.Request.RequestURI))
	})

	// 静态文件处理
	c.R.Use(static.Serve("/", static.LocalFile("./static", true)))

	// 超时中间件
	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(time.Duration(c.TimeoutDuration)*time.Second, appo.NewServiceUnavailable()))
	}

	ag := c.R.Group("api/v1/ygg")
	ag.POST("/user/register", h.Register)
	ag.POST("/login/passwd", h.Login)
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("配置session错误: %v\n", err.Error())
	}
}
