package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("!!!token has expired")
	ErrTokenInvalid = errors.New("!!!token is invalid")
)

// Payload data of token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Validates token payload (only used for PASETO)
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrTokenExpired
	}
	return nil
}
