package middleware

import (
	"github.com/bimaputraas/rest-api/internal/usecase"
	pkgstrings "github.com/bimaputraas/rest-api/pkg/strings"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type (
	Middleware struct {
		usecase *usecase.Usecase
	}
)

func New(usecase *usecase.Usecase) *Middleware {
	return &Middleware{
		usecase: usecase,
	}
}

func (m *Middleware) Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	})
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Replace(ctx.GetHeader("Authorization"), "Bearer ", "", 1)
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": pkgstrings.Capitalize("Unauthorized"),
				"details": "No token",
				"code":    401,
				"status":  "FAILED",
			})
			return
		}

		userId, err := m.usecase.Auth(ctx.Request.Context(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"details": pkgstrings.Capitalize(err.Error()),
				"code":    401,
				"status":  "FAILED",
			})
			return
		}
		ctx.Set("user_id", userId)
		ctx.Next()
	}
}
