package router

import (
	"github.com/nimil-jp/gin-utils/context"
	"github.com/nimil-jp/gin-utils/errors"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	"gorm.io/gorm"
)

type Router struct {
	g     *gin.RouterGroup
	getDB func() *gorm.DB
}

func New(engine *gin.Engine, getDB func() *gorm.DB) *Router {
	return &Router{g: engine.Group(""), getDB: getDB}
}

func (r *Router) Group(relativePath string, handlers []gin.HandlerFunc, fn func(r *Router)) {
	if handlers == nil {
		handlers = []gin.HandlerFunc{}
	}
	fn(&Router{g: r.g.Group(relativePath, handlers...), getDB: r.getDB})
}

type HandlerFunc func(ctx context.Context, c *gin.Context) error

func (r *Router) Get(relativePath string, handlerFunc HandlerFunc) {
	r.g.GET(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Post(relativePath string, handlerFunc HandlerFunc) {
	r.g.POST(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Put(relativePath string, handlerFunc HandlerFunc) {
	r.g.PUT(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Patch(relativePath string, handlerFunc HandlerFunc) {
	r.g.PATCH(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Delete(relativePath string, handlerFunc HandlerFunc) {
	r.g.DELETE(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Options(relativePath string, handlerFunc HandlerFunc) {
	r.g.OPTIONS(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Head(relativePath string, handlerFunc HandlerFunc) {
	r.g.HEAD(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) wrapperFunc(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context

		if userID, ok := jwt.GetClaim(c, "user_id"); ok {
			ctx = context.New(c.GetHeader("X-Request-Id"), uint(userID.(float64)), r.getDB)
		} else {
			ctx = context.New(c.GetHeader("X-Request-Id"), 0, r.getDB)
		}

		c.Writer.Header().Add("X-Request-Id", ctx.RequestID())

		err := handlerFunc(ctx, c)

		if err != nil {
			switch v := err.(type) {
			case *errors.Error:
				v.Response().Do(c, ctx.RequestID())
			default:
				errors.NewUnexpected(v).Response().Do(c, ctx.RequestID())
			}

			_ = c.Error(err)
		}
	}
}
