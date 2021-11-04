package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionOption struct {
	MaxAge time.Duration
}

func Session(name []string, secret string, option *SessionOption) gin.HandlerFunc {
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

	var maxAge = time.Hour * 24 * 365

	if option != nil {
		maxAge = option.MaxAge
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(
		sessions.Options{
			Path:     "/",
			MaxAge:   int(maxAge),
			Secure:   corsSecure,
			HttpOnly: true,
			SameSite: corsSameSite,
		},
	)
	return sessions.SessionsMany(name, store)
}
