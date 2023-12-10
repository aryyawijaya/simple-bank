package server

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules/account"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for banking services
type Server struct {
	Router *gin.Engine
}

// NewServer create Server and setup the routing
func NewServer(store mydb.Store) *Server {
	server := &Server{}
	router := gin.Default()

	// dependencies
	wrapper := wrapper.NewWrapper()

	// accounts
	am := account.NewAccountModule(store, wrapper)
	router.POST("/accounts", am.Create)
	router.GET("/accounts/:id", am.Get)
	router.GET("/accounts", am.GetAll)

	server.Router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
