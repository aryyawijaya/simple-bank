package auth

import (
	"database/sql"
	"net/http"

	"github.com/aryyawijaya/simple-bank/modules/user"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string            `json:"accessToken"`
	User        user.UserResponse `json:"user"`
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

	accessToken, err := am.token.CreateToken(
		req.Username,
		am.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, am.wrapper.ErrResp(err))
		return
	}

	resp := LoginResponse{
		AccessToken: accessToken,
		User:        *user.NewUserResp(&currUser),
	}

	ctx.JSON(http.StatusOK, resp)
}
