package simplemail

import (
	"encoding"
	"github.com/emersion/go-message/mail"
)

type FromAddress mail.Address

var _ encoding.TextUnmarshaler = &FromAddress{}

func (f FromAddress) String() string {
	return f.ToMailAddress().String()
}

func (f *FromAddress) UnmarshalText(b []byte) error {
	address, err := mail.ParseAddress(string(b))
	if err != nil {
		return err
	}
	*f = FromAddress(*address)
	return nil
}

func (f *FromAddress) ToMailAddress() *mail.Address {
	return (*mail.Address)(f)
}
