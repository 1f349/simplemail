package simplemail

import (
	"encoding"
	"encoding/json"
	"github.com/emersion/go-message/mail"
)

type FromAddress struct {
	*mail.Address
}

var _ encoding.TextUnmarshaler = &FromAddress{}

func (f *FromAddress) UnmarshalText(b []byte) error {
	var a string
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	address, err := mail.ParseAddress(a)
	if err != nil {
		return err
	}
	f.Address = address
	return nil
}
