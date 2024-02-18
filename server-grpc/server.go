package servergrpc

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	authdelivery "github.com/aryyawijaya/simple-bank/modules/auth/delivery"
	"github.com/aryyawijaya/simple-bank/modules/auth/password"
	authrepo "github.com/aryyawijaya/simple-bank/modules/auth/repo"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/modules/auth/token/paseto"
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	userdelivery "github.com/aryyawijaya/simple-bank/modules/user/delivery"
	userrepo "github.com/aryyawijaya/simple-bank/modules/user/repo"
	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
	"github.com/aryyawijaya/simple-bank/util"
)

// Server serves HTTP request for banking services
type Server struct {
	store      mydb.Store
	Config     *util.Config
	TokenMaker token.Maker
	*userdelivery.UserGRPC
	*authdelivery.AuthGRPC
}

// NewServer create Server and setup the routing
func NewServer(store mydb.Store, config *util.Config) (*Server, error) {
	server := &Server{
		store:  store,
		Config: config,
	}

	// other dependencies
	passHelper := password.NewPassHelper()
	paseto, err := paseto.NewPasetoMaker(server.Config.TokenSymetricKey)
	if err != nil {
		return nil, err
	}

	// set token maker
	server.TokenMaker = paseto

	// user
	userRepo := userrepo.NewUserRepo(server.store)
	userUseCase := userusecase.NewUserUseCase(passHelper, userRepo)
	userGRPC := userdelivery.NewUserGRPC(userUseCase)
	server.UserGRPC = userGRPC

	// auth
	authRepo := authrepo.NewAuthRepo(server.store)
	authUseCase := authusecase.NewAuthUseCase(authRepo, passHelper, server.TokenMaker, server.Config)
	authGRPC := authdelivery.NewAuthGRPC(authUseCase)
	server.AuthGRPC = authGRPC

	return server, nil
}
