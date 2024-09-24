package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) GetUserTransactions(ctx *gin.Context) {
	var (
		userId = ctx.MustGet("user_id").(uint)
	)

	data, err := c.usecase.GetAllUserTransactions(ctx.Request.Context(), userId)
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
