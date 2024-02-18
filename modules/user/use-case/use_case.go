package userusecase

import (
	"context"

	"github.com/aryyawijaya/simple-bank/entity"
)

type UserUseCase struct {
	passHelper PassHelper
	repo       Repo
}

func NewUserUseCase(passHelper PassHelper, repo Repo) UseCase {
	return &UserUseCase{
		passHelper: passHelper,
		repo:       repo,
	}
}

func (u *UserUseCase) Create(ctx context.Context, dto *CreateUserDto) (*entity.User, error) {
	hashedPass, err := u.passHelper.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:       dto.Username,
		HashedPassword: hashedPass,
		FullName:       dto.FullName,
		Email:          dto.Email,
	}
	createdUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
