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

	ag.Use(middleware.AuthMiddleware())
	ag.GET("/user/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result": gin.H{
				"id":               "150000200503064268",
				"name":             "尸佼",
				"avatar":           "http://dummyimage.com/64x64/FF6600",
				"email":            "drremain@crew4dance.com",
				"job":              "frontend",
				"jobName":          "前端开发工程师",
				"organization":     "Frontend",
				"organizationName": "前端组",
				"location":         "hangzhou",
				"locationName":     "杭州",
				"introduction":     "三算它白精准资影过南再战入。",
				"contact":          "13512344321",
				"lastLoginTime":    337922163112,
				"hiredType":        998071869408,
				"permissions":      "*",
			},
		})
	})
}

func setUserSession(c *gin.Context, token string) {
	session := sessions.Default(c)
	session.Set("access_token", token)
	if err := session.Save(); err != nil {
		log.Printf("配置session错误: %v\n", err.Error())
	}
}
