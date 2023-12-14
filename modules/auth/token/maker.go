package token

import "time"

// Maker interface for managing tokens
type Maker interface {
	// Creates new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// Verify token valid or not
	VerifyToken(token string) (*Payload, error)
}
