package servergrpc

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/pb"
	"github.com/aryyawijaya/simple-bank/util"
)

// Server serves HTTP request for banking services
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      mydb.Store
	Config     *util.Config
	TokenMaker token.Maker
}

// NewServer create Server and setup the routing
func NewServer(store mydb.Store, config *util.Config) (*Server, error) {
	server := &Server{
		store:  store,
		Config: config,
	}

	return server, nil
}
