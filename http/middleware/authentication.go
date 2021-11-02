package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
)

func Authentication(realm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			session := sessions.Default(c)
			token := session.Get("token")
			if token, ok := token.(string); ok {
				c.Request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
			}
		}

		jwt.Verify(realm)(c)
	}
}
