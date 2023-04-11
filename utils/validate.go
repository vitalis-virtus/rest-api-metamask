package utils

import (
	"regexp"

	"github.com/vitalis-virtus/rest-api-metamask/utils/fail"

	"github.com/vitalis-virtus/rest-api-metamask/model"
)

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
var nonceRegex *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)

func Validate(s model.SignInPayload) error {
	if !hexRegex.MatchString(s.Address) {
		return fail.ErrInvalidAddress
	}
	if !nonceRegex.MatchString(s.Nonce) {
		return fail.ErrInvalidNonce
	}
	if len(s.Sig) == 0 {
		return fail.ErrMissingSig
	}
	return nil
}
