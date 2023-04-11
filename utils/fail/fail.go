package fail

import "errors"

var (
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidNonce   = errors.New("invalid nonce")
	ErrMissingSig     = errors.New("signature is missing")
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotExists  = errors.New("user does not exist")
	ErrAuthError      = errors.New("authentication error")
)
