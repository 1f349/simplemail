package simplemail

import (
	"encoding/json"
	"net/mail"
)

type FromAddress struct {
	*mail.Address
}

var _ json.Unmarshaler = &FromAddress{}

func (f *FromAddress) UnmarshalJSON(b []byte) error {
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
