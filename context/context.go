package context

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nimil-jp/gin-utils/errors"
	"gorm.io/gorm"
)

type Context interface {
	RequestID() string
	Authenticated() bool
	UID() uint
	FirebaseUID() string

	Validate(request interface{}) (invalid bool)
	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error

	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
}

type ctx struct {
	id          string
	verr        *errors.Error
	getDB       func() *gorm.DB
	db          *gorm.DB
	uid         uint
	firebaseUID string
}

func New(appName string, c *gin.Context, getDB func() *gorm.DB) Context {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	firebaseUID := ""
	firebaseUIDInterface, ok := c.Get("firebase_uid")
	if ok {
		firebaseUID = firebaseUIDInterface.(string)
	}

	var uid uint
	claimsInterface, ok := c.Get("claims")
	if ok {
		if uidInterface, ok := claimsInterface.(map[string]interface{})[fmt.Sprintf("%s_id", appName)]; ok {
			uid = uint(uidInterface.(float64))
		}
	}

	return &ctx{
		id:          requestID,
		verr:        errors.NewValidation(),
		getDB:       getDB,
		uid:         uid,
		firebaseUID: firebaseUID,
	}
}

func (c ctx) RequestID() string {
	return c.id
}

func (c ctx) Authenticated() bool {
	return c.firebaseUID != ""
}

func (c ctx) UID() uint {
	return c.uid
}

func (c ctx) FirebaseUID() string {
	return c.firebaseUID
}
