package controller

import (
	pkgstrings "github.com/bimaputraas/rest-api/pkg/strings"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) TopUp(ctx *gin.Context) {
	var (
		payload struct {
			Amount float64 `json:"amount"`
		}
		userId = ctx.MustGet("user_id").(uint)
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

	data, err := c.usecase.TopUp(ctx.Request.Context(), userId, payload.Amount)
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
