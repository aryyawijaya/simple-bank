package wrapper

import (
	"net/http"

	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

type Wrapper struct{}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

func (w *Wrapper) ErrResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ErrResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case entity.ErrUnique:
		return http.StatusBadRequest

	case entity.ErrNotFound:
		return http.StatusNotFound

	case entity.ErrTokenExpired,
		entity.ErrSessionBlocked,
		entity.ErrSessionExpired:
		return http.StatusForbidden

	case entity.ErrTokenInvalid,
		entity.ErrSessionInvalid,
		entity.ErrPasswordInvalid:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}

func GetCodesGRPC(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	switch err {
	case entity.ErrUnique:
		return codes.AlreadyExists

	case entity.ErrNotFound:
		return codes.NotFound

	case entity.ErrTokenExpired,
		entity.ErrSessionBlocked,
		entity.ErrSessionExpired:
		return codes.PermissionDenied

	case entity.ErrTokenInvalid,
		entity.ErrSessionInvalid,
		entity.ErrPasswordInvalid:
		return codes.Unauthenticated

	default:
		return codes.Internal
	}
}
