package authdelivery

import (
	"context"

	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	"github.com/aryyawijaya/simple-bank/pb"
	utilgrpc "github.com/aryyawijaya/simple-bank/util/grpc"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"google.golang.org/grpc/status"
)

type AuthGRPC struct {
	pb.UnimplementedSimpleBankServer
	authUseCase authusecase.UseCase
}

func NewAuthGRPC(authUseCase authusecase.UseCase) *AuthGRPC {
	return &AuthGRPC{
		authUseCase: authUseCase,
	}
}

func (u *AuthGRPC) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	mtdt := utilgrpc.ExtractMetadata(ctx)
	
	dto := &authusecase.LoginDto{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		UserAgent: mtdt.UserAgent,
		ClientIP:  mtdt.ClientIP,
	}
	loggedUser, err := u.authUseCase.Login(ctx, dto)
	if err != nil {
		return nil, status.Errorf(wrapper.GetCodesGRPC(err), err.Error())
	}

	resp := newPbLoginResponse(loggedUser)

	return resp, nil
}
