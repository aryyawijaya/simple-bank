package userdelivery

import (
	"net/http"

	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
	"github.com/aryyawijaya/simple-bank/util/adapter"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	userUseCase userusecase.UseCase
}

func NewUserHttp(router *gin.Engine, userUseCase userusecase.UseCase) {
	userHttp := &UserHttp{
		userUseCase: userUseCase,
	}

	router.POST("/users", userHttp.Create)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (u *UserHttp) Create(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapper.ErrResp(err))
		return
	}

	dto := &userusecase.CreateUserDto{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
	}
	createdUser, err := u.userUseCase.Create(ctx, dto)
	if err != nil {
		ctx.JSON(wrapper.GetStatusCode(err), wrapper.ErrResp(err))
		return
	}

	resp := adapter.NewUserResp(createdUser)

	ctx.JSON(http.StatusCreated, resp)
}
