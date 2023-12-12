package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthHelper struct{}

func NewAuthHelper() *AuthHelper {
	return &AuthHelper{}
}

func (ah *AuthHelper) HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPass), nil
}

func (ah *AuthHelper) CheckPassword(hashedPass string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
