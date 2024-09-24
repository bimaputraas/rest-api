package controller

import (
	"github.com/bimaputraas/rest-api/internal/usecase"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkgstrings "github.com/bimaputraas/rest-api/pkg/strings"
	"github.com/gin-gonic/gin"
)

type (
	Controller struct {
		usecase *usecase.Usecase
	}
)

func New(usecase *usecase.Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

func errUsecaseHandler(err error) (int, gin.H) {
	errUsecase, ok := pkgerrors.ParseError(err)
	if !ok {
		return 500, gin.H{
			"message": "Internal Server Error",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    500,
		}
	}
	switch errUsecase.Code() {
	case pkgerrors.ErrInvalidArgument:
		return 400, gin.H{
			"message": "Bad Request",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    400,
		}
	case pkgerrors.ErrNotFound:
		return 404, gin.H{
			"message": "Not Found",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    404,
		}
	case pkgerrors.ErrIllegal:
		return 401, gin.H{
			"message": "Unauthorized",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    401,
		}
	default:
		return 500, gin.H{
			"message": "Internal Server Error",
			"details": pkgstrings.Capitalize(err.Error()),
			"code":    500,
		}
	}
}
