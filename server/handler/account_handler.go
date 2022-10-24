package handler

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go-remix/model/apperrors"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) GetCurrent(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	user, err := h.userService.Get(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

type editReq struct {
	Username string `form:"username"`
	Email    string `form:"email"`
}

func (r editReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
	)
}

func (r *editReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
}

func (h *Handler) Edit(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	var req editReq
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	authUser, err := h.userService.Get(userId)

	if err != nil {
		e := apperrors.NewAuthorization(apperrors.InvalidSession)
		c.JSON(e.Status(), gin.H{
			"errror": e,
		})
		return
	}

	authUser.Username = req.Username
	if authUser.Email != req.Email {
		inUse := h.userService.IsEmailAlreadyInUse(req.Email)
		if inUse {
			toFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		authUser.Email = req.Email
	}

	err = h.userService.UpdateAccount(authUser)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, authUser)
}

type changeRequest struct {
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

func (r changeRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CurrentPassword, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.ConfirmNewPassword, validation.Required, validation.Length(6, 150)),
	)
}

func (r *changeRequest) sanitize() {
	r.CurrentPassword = strings.TrimSpace(r.CurrentPassword)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.ConfirmNewPassword = strings.TrimSpace(r.ConfirmNewPassword)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	var req changeRequest

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	if req.NewPassword != req.ConfirmNewPassword {
		toFieldErrorResponse(c, "password", apperrors.PasswordsDoNotMatch)
		return
	}

	authUser, err := h.userService.Get(userId)

	if err != nil {
		e := apperrors.NewAuthorization(apperrors.InvalidSession)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	err = h.userService.ChangePassword(req.CurrentPassword, req.NewPassword, authUser)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
