package simplemail

import (
	"encoding"
	"github.com/emersion/go-message/mail"
)

type FromAddress struct {
	*mail.Address
}

var _ encoding.TextUnmarshaler = &FromAddress{}

func (f *FromAddress) UnmarshalText(b []byte) error {
	address, err := mail.ParseAddress(string(b))
	if err != nil {
		return err
	}
	f.Address = address
	return nil
}
