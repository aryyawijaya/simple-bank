package server

import (
	"github.com/aryyawijaya/simple-bank/middlewares"
	"github.com/aryyawijaya/simple-bank/modules/account"
	"github.com/aryyawijaya/simple-bank/modules/auth"
	"github.com/aryyawijaya/simple-bank/modules/auth/password"
	"github.com/aryyawijaya/simple-bank/modules/auth/token/paseto"
	"github.com/aryyawijaya/simple-bank/modules/transfer"
	"github.com/aryyawijaya/simple-bank/modules/user"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRouter() error {
	router := gin.Default()

	// dependencies
	wrapper := wrapper.NewWrapper()
	passHelper := password.NewPassHelper()
	paseto, err := paseto.NewPasetoMaker(s.Config.TokenSymetricKey)
	if err != nil {
		return err
	}
	middleware := middlewares.NewMiddleware(wrapper)

	// set token maker
	s.TokenMaker = paseto

	// router group
	authRouter := router.Group("/").Use(middleware.AuthMiddleware(s.TokenMaker))

	// users
	um := user.NewUserModule(s.store, wrapper, passHelper)
	router.POST("/users", um.Create)

	// accounts
	am := account.NewAccountModule(s.store, wrapper)
	authRouter.POST("/accounts", am.Create)
	authRouter.GET("/accounts/:id", am.Get)
	authRouter.GET("/accounts", am.GetAll)

	// transfer
	tm := transfer.NewTransferModule(s.store, wrapper)
	authRouter.POST("/transfers", tm.CreateTransfer)

	// auth
	authModule := auth.NewAuthModule(s.Config, wrapper, s.store, passHelper, s.TokenMaker)
	router.POST("/login", authModule.Login)

	s.Router = router

	return nil
}
