package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"go-remix/appo"
	"go-remix/config"
	"go-remix/middleware"
	"go-remix/model"
)

type Handler struct {
	userService  model.UserService
	MaxBodyBytes int64
}

func InjectRouter(c *gin.Engine, cfg config.Config) {
	_ = &Handler{
		MaxBodyBytes: cfg.MaxBodyBytes,
	}

	// 不存在路由处理
	c.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, appo.NewNotFound("api", c.Request.RequestURI))
	})

	// 静态文件处理
	c.Use(static.Serve("/", static.LocalFile("./static", true)))

	// 超时中间件
	if gin.Mode() != gin.TestMode {
		c.Use(middleware.Timeout(time.Duration(cfg.HandlerTimeOut)*time.Second, appo.NewServiceUnavailable()))
	}

	//ag := c.Group("api/account")
	//ag.POST("/register", h.Register)
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("配置session错误: %v\n", err.Error())
	}
}
