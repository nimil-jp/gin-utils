package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-Id"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}
