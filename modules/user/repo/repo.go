package userrepo

import (
	"context"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/entity"
	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
	"github.com/aryyawijaya/simple-bank/util/adapter"
	"github.com/lib/pq"
)

type UserRepo struct {
	store mydb.Store
}

func NewUserRepo(store mydb.Store) userusecase.Repo {
	return &UserRepo{
		store: store,
	}
}

func (u *UserRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	arg := mydb.CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	createdUser, err := u.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, entity.ErrUnique
			}
		}
		return nil, err
	}

	userEntity := adapter.UserSqlcToEntity(&createdUser)
	return userEntity, nil
}
