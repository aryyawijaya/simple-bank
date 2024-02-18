package authdelivery

import (
	"net/http"
	"time"

	"github.com/aryyawijaya/simple-bank/entity"
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	"github.com/aryyawijaya/simple-bank/util/adapter"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHttp struct {
	authUseCase authusecase.UseCase
}

func NewAuthHttp(router *gin.Engine, authUseCase authusecase.UseCase) {
	authHttp := &AuthHttp{
		authUseCase: authUseCase,
	}

	router.POST("/login", authHttp.Login)
	router.POST("/tokens/renew-access", authHttp.RenewAccessToken)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	SessionID             uuid.UUID            `json:"sessionId"`
	AccessToken           string               `json:"accessToken"`
	AccessTokenExpiresAt  time.Time            `json:"accessTokenExpiresAt"`
	RefreshToken          string               `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time            `json:"refreshTokenExpiresAt"`
	User                  *entity.UserResponse `json:"user"`
}

func (a *AuthHttp) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapper.ErrResp(err))
		return
	}

	dto := &authusecase.LoginDto{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
	}
	loggedUser, err := a.authUseCase.Login(ctx, dto)
	if err != nil {
		ctx.JSON(wrapper.GetStatusCode(err), wrapper.ErrResp(err))
		return
	}

	resp := &LoginResponse{
		SessionID:             loggedUser.SessionID,
		AccessToken:           loggedUser.AccessToken,
		AccessTokenExpiresAt:  loggedUser.AccessTokenExpiresAt,
		RefreshToken:          loggedUser.RefreshToken,
		RefreshTokenExpiresAt: loggedUser.RefreshTokenExpiresAt,
		User:                  adapter.NewUserResp(loggedUser.User),
	}

	ctx.JSON(http.StatusOK, resp)
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func (a *AuthHttp) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapper.ErrResp(err))
		return
	}

	renewedAccessToken, err := a.authUseCase.RenewAccessToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(wrapper.GetStatusCode(err), wrapper.ErrResp(err))
		return
	}

	resp := &RenewAccessTokenResponse{
		AccessToken:          renewedAccessToken.AccessToken,
		AccessTokenExpiresAt: renewedAccessToken.AccessTokenExpiresAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
