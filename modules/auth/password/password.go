package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PassHelper struct{}

func NewPassHelper() *PassHelper {
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
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
