package api

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for banking services
type Server struct {
	store  *mydb.Store
	router *gin.Engine
}

// NewServer create Server and setup the routing
func NewServer(store *mydb.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
