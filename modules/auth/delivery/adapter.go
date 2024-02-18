package authdelivery

import (
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	"github.com/aryyawijaya/simple-bank/pb"
	"github.com/aryyawijaya/simple-bank/util/adapter"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newPbLoginResponse(loginResponse *authusecase.LoginResponse) *pb.LoginResponse {
	return &pb.LoginResponse{
		SessionId:             loginResponse.SessionID.String(),
		AccessToken:           loginResponse.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(loginResponse.AccessTokenExpiresAt),
		RefreshToken:          loginResponse.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(loginResponse.RefreshTokenExpiresAt),
		User:                  adapter.UserEntityToPb(loginResponse.User),
	}
}
