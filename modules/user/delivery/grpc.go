package userdelivery

import (
	"context"

	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
	"github.com/aryyawijaya/simple-bank/pb"
	utilgrpc "github.com/aryyawijaya/simple-bank/util/grpc"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"google.golang.org/grpc/status"
)

type UserGRPC struct {
	userUseCase userusecase.UseCase
}

func NewUserGRPC(userUseCase userusecase.UseCase) *UserGRPC {
	return &UserGRPC{
		userUseCase: userUseCase,
	}
}

func (u *UserGRPC) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// validate request
	violations := validateCreateUserRequest(req)
	if len(violations) != 0 {
		return nil, utilgrpc.InvalidArgumentErr(violations)
	}

	dto := &userusecase.CreateUserDto{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}
	createdUser, err := u.userUseCase.Create(ctx, dto)
	if err != nil {
		return nil, status.Errorf(wrapper.GetCodesGRPC(err), err.Error())
	}

	resp := newPbCreateUserResponse(createdUser)

	return resp, nil
}
