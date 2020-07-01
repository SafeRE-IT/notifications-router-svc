package types

import (
	"errors"

	"github.com/spf13/cast"
	"gitlab.com/tokend/go/strkey"
)

var (
	ErrAccountIDInvalid = errors.New("address is invalid")
)

type AccountID string

func (a AccountID) Validate() error {
	_, err := strkey.Decode(strkey.VersionByteAccountID, string(a))
	if err != nil {
		return ErrAccountIDInvalid
	}
	return nil
}

func (a AccountID) String() string {
	return string(a)
}

var IsAccountID = &isAccountID{}

type isAccountID struct{}

func (ia *isAccountID) Validate(value interface{}) error {
	a, err := cast.ToStringE(value)
	if err != nil {
		return err
	}
	address := AccountID(a)
	return address.Validate()
}
