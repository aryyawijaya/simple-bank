package middlewares

import "github.com/gin-gonic/gin"

type Wrapper interface {
	ErrResp(err error) gin.H
}

type Middleware struct {
	wrapper Wrapper
}

func NewMiddleware(wrapper Wrapper) *Middleware {
	return &Middleware{
		wrapper: wrapper,
	}
}
