package user

import (
	"context"
	"net/http"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Store interface {
	CreateUser(ctx context.Context, arg mydb.CreateUserParams) (mydb.User, error)
}

type PassHelper interface {
	HashPassword(password string) (string, error)
}

type UserModule struct {
	store      Store
	wrapper    modules.Wrapper
	passHelper PassHelper
}

func NewUserModule(store Store, wrapper modules.Wrapper, passHelper PassHelper) *UserModule {
	return &UserModule{
		store:      store,
		wrapper:    wrapper,
		passHelper: passHelper,
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (um *UserModule) Create(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, um.wrapper.ErrResp(err))
		return
	}

	hashedPass, err := um.passHelper.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, um.wrapper.ErrResp(err))
		return
	}

	arg := mydb.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPass,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	createdUser, err := um.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, um.wrapper.ErrResp(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, um.wrapper.ErrResp(err))
		return
	}

	resp := NewUserResp(&createdUser)

	ctx.JSON(http.StatusCreated, resp)
}
