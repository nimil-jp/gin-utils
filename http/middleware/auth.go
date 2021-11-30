package middleware

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type FirebaseIDTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

var firebaseAuthClient FirebaseIDTokenVerifier

func FirebaseSetup(client FirebaseIDTokenVerifier) {
	firebaseAuthClient = client
}

func FirebaseAuth(must bool) gin.HandlerFunc {
	if firebaseAuthClient == nil {
		panic("FirebaseAuthの前にFirebaseSetupを実行してください")
	}

	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		nextAction := func() {
			if must {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
		}

		if len(authorizationHeader) <= 7 {
			nextAction()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		token, err := firebaseAuthClient.VerifyIDToken(ctx, authorizationHeader[7:])
		if err != nil {
			nextAction()
			return
		}
		c.Set("firebase_uid", token.UID)
		c.Set("claims", token.Claims)
		c.Next()
	}
}
