package userusecase

import (
	"context"

	"github.com/aryyawijaya/simple-bank/entity"
)

type UseCase interface {
	Create(ctx context.Context, dto *CreateUserDto) (*entity.User, error)
}

type Repo interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
}

type PassHelper interface {
	HashPassword(password string) (string, error)
}
