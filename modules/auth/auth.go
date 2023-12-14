package auth

import (
	"context"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/aryyawijaya/simple-bank/util"
)

type Store interface {
	GetUser(ctx context.Context, username string) (mydb.User, error)
}

type PassHelper interface {
	CheckPassword(hashedPass string, pass string) error
}

type Token interface {
	CreateToken(username string, duration time.Duration) (string, error)
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
