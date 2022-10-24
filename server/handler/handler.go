package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-remix/handler/middleware"
	"go-remix/model"
	"go-remix/model/apperrors"
	"log"
	"net/http"

	"time"
)

type Handler struct {
	userService  model.UserService
	MaxBodyBytes int64
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	h := &Handler{
		userService:  c.UserService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found.",
		})
	})

	c.R.Use(static.Serve("/", static.LocalFile("./static", true)))

	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	}

	c.R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ag := c.R.Group("api/account")

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)
	ag.POST("/logout", h.Logout)
	ag.POST("/forgot-password", h.ForgotPassword)
	ag.POST("/reset-password", h.ResetPassword)

	ag.Use(middleware.AuthUser())
	ag.GET("", h.GetCurrent)
	ag.PUT("", h.Edit)
	ag.PUT("/change-password", h.ChangePassword)
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("error setting the session: %v\n", err.Error())
	}
}

func toFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []model.FieldError{
			{
				Field: field, Message: message,
			},
		},
	})
}
