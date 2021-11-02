package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session(name, secret string) gin.HandlerFunc {
	var (
		corsSecure   bool
		corsSameSite http.SameSite
	)

	switch gin.Mode() {
	case gin.ReleaseMode:
		corsSecure = true
		corsSameSite = http.SameSiteStrictMode
	case gin.DebugMode:
		corsSecure = false
		corsSameSite = http.SameSiteLaxMode
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(
		sessions.Options{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365,
			Secure:   corsSecure,
			HttpOnly: true,
			SameSite: corsSameSite,
		},
	)
	return sessions.Sessions(name, store)
}
