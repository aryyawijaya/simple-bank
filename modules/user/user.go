package user

import (
	"context"
	"net/http"
	"time"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/modules"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Store interface {
	CreateUser(ctx context.Context, arg mydb.CreateUserParams) (mydb.User, error)
}

type AuthHelper interface {
	HashPassword(password string) (string, error)
}

type UserModule struct {
	store      Store
	wrapper    modules.Wrapper
	authHelper AuthHelper
}

func NewUserModule(store Store, wrapper modules.Wrapper, authHelper AuthHelper) *UserModule {
	return &UserModule{
		store:      store,
		wrapper:    wrapper,
		authHelper: authHelper,
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (um *UserModule) Create(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, um.wrapper.ErrResp(err))
		return
	}

	hashedPass, err := um.authHelper.HashPassword(req.Password)
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

	resp := createUserResponse{
		Username: createdUser.Username,
		FullName: createdUser.FullName,
		Email: createdUser.Email,
		PasswordChangedAt: createdUser.PasswordChangedAt,
		CreatedAt: createdUser.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, resp)
}
