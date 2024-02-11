package authusecase

import (
	"context"
	"time"

	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/aryyawijaya/simple-bank/util"
)

type AuthUseCase struct {
	repo       Repo
	passHelper PassHelper
	token      Token
	config     *util.Config
}

func NewAuthUseCase(repo Repo, passHelper PassHelper, token Token, config *util.Config) UseCase {
	return &AuthUseCase{
		repo:       repo,
		passHelper: passHelper,
		token:      token,
		config:     config,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, dto *LoginDto) (*LoginResponse, error) {
	currUser, err := a.repo.GetUser(ctx, dto.Username)
	if err != nil {
		return nil, err
	}

	err = a.passHelper.CheckPassword(currUser.HashedPassword, dto.Password)
	if err != nil {
		return nil, err
	}

	// access token
	accessToken, accessPayload, err := a.token.CreateToken(
		dto.Username,
		a.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshToken, refreshPayload, err := a.token.CreateToken(
		dto.Username,
		a.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// session
	session, err := a.repo.CreateSession(ctx, &CreateSessionDto{
		ID:           refreshPayload.ID,
		Username:     currUser.Username,
		RefreshToken: refreshToken,
		UserAgent:    dto.UserAgent,
		ClientIp:     dto.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	logged := &LoginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  currUser,
	}

	return logged, nil
}

func (a *AuthUseCase) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenResponse, error) {
	// validate refresh token by Token Implementor (Paseto or JWT)
	refreshPayload, err := a.token.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// validate refresh token by record
	currSession, err := a.repo.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		return nil, err
	}
	if currSession.IsBlocked {
		return nil, entity.ErrSessionBlocked
	}
	if currSession.Username != refreshPayload.Username {
		return nil, entity.ErrSessionInvalid
	}
	if currSession.RefreshToken != refreshToken {
		return nil, entity.ErrSessionInvalid
	}
	// check expires time 1 more time
	if time.Now().After(currSession.ExpiresAt) {
		return nil, entity.ErrSessionExpired
	}

	// access token
	accessToken, accessPayload, err := a.token.CreateToken(
		refreshPayload.Username,
		a.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	renewedAccessToken := &RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return renewedAccessToken, nil
}
