package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JWTMaker implements Maker interface from token package
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (token.Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf(
			"invalid key size: must be at least %d characters",
			minSecretKeySize,
		)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

type JWTPayload struct {
	token.Payload
	jwt.RegisteredClaims
}

func (jm *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := token.NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	claims := JWTPayload{
		*payload,
		jwt.RegisteredClaims{
			ID:        payload.ID.String(),
			IssuedAt:  &jwt.NumericDate{Time: payload.IssuedAt},
			ExpiresAt: &jwt.NumericDate{Time: payload.ExpiredAt},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(jm.secretKey))
}

func (jm *JWTMaker) VerifyToken(t string) (*token.Payload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, token.ErrTokenInvalid
		}

		return []byte(jm.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(t, &JWTPayload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, token.ErrTokenExpired
		}
		return nil, token.ErrTokenInvalid
	}

	claims, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		// err = errors.New("unknown claims type, cannot proceed")
		return nil, token.ErrTokenInvalid
	}

	return &claims.Payload, nil
}
