package password

import (
	"fmt"

	"github.com/aryyawijaya/simple-bank/entity"
	authusecase "github.com/aryyawijaya/simple-bank/modules/auth/use-case"
	userusecase "github.com/aryyawijaya/simple-bank/modules/user/use-case"
	"golang.org/x/crypto/bcrypt"
)

type IPassHelper interface {
	userusecase.PassHelper
	authusecase.PassHelper
}

type PassHelper struct{}

func NewPassHelper() IPassHelper {
	return &PassHelper{}
}

func (ph *PassHelper) HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPass), nil
}

func (ph *PassHelper) CheckPassword(hashedPass string, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return entity.ErrPasswordInvalid
		}
		return err
	}

	return nil
}
