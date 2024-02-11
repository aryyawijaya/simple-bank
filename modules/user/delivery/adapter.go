package userdelivery

import (
	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/aryyawijaya/simple-bank/pb"
	"github.com/aryyawijaya/simple-bank/util/adapter"
)

func newPbCreateUserResponse(user *entity.User) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		User: adapter.UserEntityToPb(user),
	}
}
