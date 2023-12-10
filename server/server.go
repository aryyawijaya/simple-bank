package server

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules/account"
	"github.com/aryyawijaya/simple-bank/modules/transfer"
	cutomvalidator "github.com/aryyawijaya/simple-bank/util/cutom-validator"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking services
type Server struct {
	Router *gin.Engine
}

// NewServer create Server and setup the routing
func NewServer(store mydb.Store) *Server {
	server := &Server{}
	router := gin.Default()

	// custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", cutomvalidator.ValidCurrency)
	}

	// dependencies
	wrapper := wrapper.NewWrapper()

	// accounts
	am := account.NewAccountModule(store, wrapper)
	router.POST("/accounts", am.Create)
	router.GET("/accounts/:id", am.Get)
	router.GET("/accounts", am.GetAll)

	// transfer
	tm := transfer.NewTransferModule(store, wrapper)
	router.POST("/transfers", tm.CreateTransfer)

	server.Router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
