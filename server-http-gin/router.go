package server

import (
	"github.com/aryyawijaya/simple-bank/middlewares"
	"github.com/aryyawijaya/simple-bank/modules/account"
	authdelivery "github.com/aryyawijaya/simple-bank/modules/auth/delivery"
	"github.com/aryyawijaya/simple-bank/modules/auth/password"
	authrepo "github.com/aryyawijaya/simple-bank/modules/auth/repo"
	"github.com/aryyawijaya/simple-bank/modules/auth/token/paseto"
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	"github.com/aryyawijaya/simple-bank/modules/transfer"
	userdelivery "github.com/aryyawijaya/simple-bank/modules/user/delivery"
	userrepo "github.com/aryyawijaya/simple-bank/modules/user/repo"
	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
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
	userRepo := userrepo.NewUserRepo(s.store)
	userUseCase := userusecase.NewUserUseCase(passHelper, userRepo)
	userdelivery.NewUserHttp(router, userUseCase)

	// accounts
	am := account.NewAccountModule(s.store, wrapper)
	authRouter.POST("/accounts", am.Create)
	authRouter.GET("/accounts/:id", am.Get)
	authRouter.GET("/accounts", am.GetAll)

	// transfer
	tm := transfer.NewTransferModule(s.store, wrapper)
	authRouter.POST("/transfers", tm.CreateTransfer)

	// auth
	authRepo := authrepo.NewAuthRepo(s.store)
	authUseCase := authusecase.NewAuthUseCase(authRepo, passHelper, s.TokenMaker, s.Config)
	authdelivery.NewAuthHttp(router, authUseCase)

	s.Router = router

	return nil
}
