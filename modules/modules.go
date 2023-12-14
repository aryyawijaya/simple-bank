package modules

import "github.com/gin-gonic/gin"

const (
	AuthorizationPayloadKey = "authorization_payload"
)

type Wrapper interface {
	ErrResp(err error) gin.H
}
