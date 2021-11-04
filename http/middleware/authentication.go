package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
)

func Authentication(jwtRealm string, session string) gin.HandlerFunc {
	if session != "" {
		return func(c *gin.Context) {
			if c.GetHeader("Authorization") == "" {
				session := sessions.DefaultMany(c, session)
				token := session.Get("token")
				if token, ok := token.(string); ok {
					c.Request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
				}
			}

			jwt.Verify(jwtRealm)(c)
		}
	} else {
		return jwt.Verify(jwtRealm)
	}
}
