package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func (am *Middleware) AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract authorization header from request
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			// client doesnt provide authorization header
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
			return
		}

		// validate authorization header format
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
			return
		}

		// validate authorization type
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", fields[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
			return
		}

		// validate access token
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
			return
		}

		// store payload in gin context
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
