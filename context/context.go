package context

import (
	"github.com/nimil-jp/gin-utils/xerrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Context interface {
	RequestID() string
	Authenticated() bool
	UserID() uint

	Validate(request interface{}) (ok bool)
	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error

	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
}

type ctx struct {
	id     string
	verr   *xerrors.Validation
	getDB  func() *gorm.DB
	db     *gorm.DB
	userID uint
}

func New(requestID string, userID uint, getDB func() *gorm.DB) Context {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return &ctx{
		id:     requestID,
		verr:   xerrors.NewValidation(),
		getDB:  getDB,
		userID: userID,
	}
}

func (c ctx) RequestID() string {
	return c.id
}

func (c ctx) Authenticated() bool {
	return c.userID != 0
}

func (c ctx) UserID() uint {
	return c.userID
}
