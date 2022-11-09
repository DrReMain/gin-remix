package handler

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go-remix/appo"
	"go-remix/model"
	"go-remix/utils"
	"log"
	"net/http"
	"strings"
)

type registerReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r registerReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	initial := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.userService.Register(initial)
	if err != nil {
		if err.Error() == appo.NewBadRequest(appo.MsgDuplicateEmail).Error() {
			c.JSON(http.StatusBadRequest, appo.NewBadRequest(appo.MsgDuplicateEmail))
			return
		}

		c.JSON(appo.Status(err), appo.NewInternal())
		return
	}

	c.JSON(http.StatusCreated, user)
}

type loginReq struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

func (r loginReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *loginReq) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(appo.Status(err), err)
		return
	}

	token, err := utils.ReleaseToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, appo.NewInternal())
		log.Printf("token generate error: %v", err)
		return
	}

	setUserSession(c, token)

	c.JSON(http.StatusOK, appo.NewSuccess(gin.H{
		"accessToken": token,
	}))
}
