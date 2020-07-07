package types

import (
	"errors"
	"regexp"
)

var (
	ErrPhoneNumberInvalid = errors.New("must be valid phone number")
)

func IsPhoneNumber(value string) bool {
	re := regexp.MustCompile(`^\+\d[1-9]\d{6,14}$`)
	return re.MatchString(value)
}

func ValidatePhoneNumber(value string) error {
	if !IsPhoneNumber(value) {
		return ErrPhoneNumberInvalid
	}
	return nil
}
