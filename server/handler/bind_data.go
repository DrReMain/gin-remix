package handler

import (
	"github.com/gin-gonic/gin"
	"go-remix/appo"
	"net/http"
)

type Request interface {
	validate() error
}

func bindData(c *gin.Context, req Request) bool {
	if err := c.ShouldBind(req); err != nil {
		c.JSON(appo.Status(err), appo.NewInternal())
		return false
	}

	if err := req.validate(); err != nil {
		c.JSON(http.StatusOK, appo.NewBadRequest(err.Error()))
		return false
	}

	return true
}
