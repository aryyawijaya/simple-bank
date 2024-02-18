package authrepo

import (
	"context"
	"database/sql"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/entity"
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	"github.com/aryyawijaya/simple-bank/util/adapter"
	"github.com/google/uuid"
)

type AuthRepo struct {
	store mydb.Store
}

func NewAuthRepo(store mydb.Store) authusecase.Repo {
	return &AuthRepo{
		store: store,
	}
}

func (a *AuthRepo) GetUser(ctx context.Context, username string) (*entity.User, error) {
	currUser, err := a.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	userEntity := adapter.UserSqlcToEntity(&currUser)

	return userEntity, nil
}

func (a *AuthRepo) CreateSession(ctx context.Context, dto *authusecase.CreateSessionDto) (*entity.Session, error) {
	session, err := a.store.CreateSession(ctx, mydb.CreateSessionParams{
		ID:           dto.ID,
		Username:     dto.Username,
		RefreshToken: dto.RefreshToken,
		UserAgent:    dto.UserAgent,
		ClientIp:     dto.ClientIp,
		IsBlocked:    dto.IsBlocked,
		ExpiresAt:    dto.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	sessionEntity := adapter.SessionSqlcToEntity(&session)

	return sessionEntity, nil
}

func (a *AuthRepo) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	currSession, err := a.store.GetSession(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	sessionEntity := adapter.SessionSqlcToEntity(&currSession)

	return sessionEntity, nil
}
