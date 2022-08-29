package authmod

import (
	"time"
)

const (
	PasswordSaltSize     = 16
	PasswordNonceSize    = 12
	PasswordNonceTimeout = 5 * time.Second

	SaltDummyString = "there is no user with this username"
)

type PasswordNonce struct {
	ID    string
	Value []byte

	Username string
	Salt     []byte
}
