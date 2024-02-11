package authusecase

import (
	"context"
	"time"

	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/google/uuid"
)

type UseCase interface {
	Login(ctx context.Context, dto *LoginDto) (*LoginResponse, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenResponse, error)
}

type Repo interface {
	GetUser(ctx context.Context, username string) (*entity.User, error)
	CreateSession(ctx context.Context, dto *CreateSessionDto) (*entity.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error)
}

type PassHelper interface {
	CheckPassword(hashedPass string, pass string) error
}

type Token interface {
	CreateToken(username string, duration time.Duration) (string, *token.Payload, error)
	VerifyToken(token string) (*token.Payload, error)
}
