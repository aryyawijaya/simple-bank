package adapter

import (
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/entity"
	"github.com/aryyawijaya/simple-bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUserResp(user *entity.User) *entity.UserResponse {
	return &entity.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func UserSqlcToEntity(user *mydb.User) *entity.User {
	return &entity.User{
		Username:          user.Username,
		HashedPassword:    user.HashedPassword,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func SessionSqlcToEntity(session *mydb.Session) *entity.Session {
	return &entity.Session{
		ID:           session.ID,
		Username:     session.Username,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func UserEntityToPb(user *entity.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
