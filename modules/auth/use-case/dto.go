package authusecase

import (
	"time"

	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/google/uuid"
)

type LoginDto struct {
	Username  string
	Password  string
	UserAgent string
	ClientIP  string
}

type LoginResponse struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  *entity.User
}

type CreateSessionDto struct {
	ID           uuid.UUID
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

type RenewAccessTokenResponse struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}
