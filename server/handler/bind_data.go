package handler

import (
	"github.com/gin-gonic/gin"
	"go-remix/model"
	"go-remix/model/apperrors"
	"net/http"
	"strings"
)

type Request interface {
	validate() error
}

func bindData(c *gin.Context, req Request) bool {
	if err := c.ShouldBind(req); err != nil {
		// TODO: 内部错误待修改
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return false
	}

	if err := req.validate(); err != nil {
		errors := strings.Split(err.Error(), ";")
		fErrors := make([]model.FieldError, 0)

		for _, e := range errors {
			split := strings.Split(e, ":")
			er := model.FieldError{
				Field:   strings.TrimSpace(split[0]),
				Message: strings.TrimSpace(split[1]),
			}
			fErrors = append(fErrors, er)
		}

		// TODO: 内部参数错误待修改
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fErrors,
		})
		return false
	}
	return true
}
