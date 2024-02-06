package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	SessionID             uuid.UUID         `json:"sessionId"`
	AccessToken           string            `json:"accessToken"`
	AccessTokenExpiresAt  time.Time         `json:"accessTokenExpiresAt"`
	RefreshToken          string            `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time         `json:"refreshTokenExpiresAt"`
	User                  user.UserResponse `json:"user"`
}

func (am *AuthModule) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, am.wrapper.ErrResp(err))
		return
	}

	currUser, err := am.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, am.wrapper.ErrResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	err = am.passHelper.CheckPassword(currUser.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}

	// access token
	accessToken, accessPayload, err := am.token.CreateToken(
		req.Username,
		am.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	// refresh token
	refreshToken, refreshPayload, err := am.token.CreateToken(
		req.Username,
		am.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	// session
	session, err := am.store.CreateSession(ctx, mydb.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     currUser.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	resp := LoginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  *user.NewUserResp(&currUser),
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

func (am *AuthModule) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, am.wrapper.ErrResp(err))
		return
	}

	// validate refresh token by Token Implementor (Paseto or JWT)
	refreshPayload, err := am.token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}

	// validate refresh token by record
	currSession, err := am.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, am.wrapper.ErrResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}
	if currSession.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}
	if currSession.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}
	if currSession.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}
	// check expires time 1 more time
	if time.Now().After(currSession.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}

	// access token
	accessToken, accessPayload, err := am.token.CreateToken(
		refreshPayload.Username,
		am.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, am.wrapper.ErrResp(err))
		return
	}

	resp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
