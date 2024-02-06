package paseto

import (
	"fmt"
	"time"

	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker implements Maker interface from token package
type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (token.Maker, error) {
	if len(symetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf(
			"invalid key size: must be at least %d characters",
			chacha20poly1305.KeySize,
		)
	}

	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}

	return maker, nil
}

func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *token.Payload, error) {
	payload, err := token.NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := pm.paseto.Encrypt(pm.symetricKey, payload, nil)

	return token, payload, err
}

func (pm *PasetoMaker) VerifyToken(t string) (*token.Payload, error) {
	payload := &token.Payload{}

	if err := pm.paseto.Decrypt(t, pm.symetricKey, payload, nil); err != nil {
		return nil, token.ErrTokenInvalid
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
