package modules

import "github.com/gin-gonic/gin"

type Wrapper interface {
	ErrResp(err error) gin.H
}
