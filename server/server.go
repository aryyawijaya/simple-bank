package server

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/aryyawijaya/simple-bank/util"
	customvalidator "github.com/aryyawijaya/simple-bank/util/cutom-validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking services
type Server struct {
	Router     *gin.Engine
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

	// custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", customvalidator.ValidCurrency)
	}

	err := server.setupRouter()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
