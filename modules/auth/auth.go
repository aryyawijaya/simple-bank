package auth

import (
	"context"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/util"
	"github.com/google/uuid"
)

type Store interface {
	GetUser(ctx context.Context, username string) (mydb.User, error)
	CreateSession(ctx context.Context, arg mydb.CreateSessionParams) (mydb.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (mydb.Session, error)
}

type PassHelper interface {
	CheckPassword(hashedPass string, pass string) error
}

type Token interface {
	CreateToken(username string, duration time.Duration) (string, *token.Payload, error)
	VerifyToken(token string) (*token.Payload, error)
}

type AuthModule struct {
	config *util.Config

	wrapper    modules.Wrapper
	store      Store
	passHelper PassHelper
	token      Token
}

func NewAuthModule(config *util.Config, wrapper modules.Wrapper, store Store, passHelper PassHelper, token Token) *AuthModule {
	return &AuthModule{
		config:     config,
		wrapper:    wrapper,
		store:      store,
		passHelper: passHelper,
		token:      token,
	}
}
