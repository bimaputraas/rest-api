package controller

import (
	"github.com/bimaputraas/rest-api/internal/model"
	"github.com/bimaputraas/rest-api/internal/usecase"
	pkgstrings "github.com/bimaputraas/rest-api/pkg/strings"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) Register(ctx *gin.Context) {
	var (
		payload = model.User{}
	)
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    400,
			"status":  "FAILED",
		})
	}

	data, err := c.usecase.Register(ctx.Request.Context(), payload)
	if err != nil {
		code, resp := errUsecaseHandler(err)
		ctx.JSON(code, resp)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Created",
		"result":  data,
		"code":    201,
		"status":  "SUCCESS",
	})
}

func (c *Controller) Login(ctx *gin.Context) {
	var (
		payload = usecase.Login{}
	)
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    400,
			"status":  "FAILED",
		})
	}

	data, err := c.usecase.Login(ctx.Request.Context(), payload)
	if err != nil {
		code, resp := errUsecaseHandler(err)
		ctx.JSON(code, resp)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
		"result":  data,
		"code":    200,
		"status":  "SUCCESS",
	})
}
