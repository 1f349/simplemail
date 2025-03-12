package simplemail

import (
	"context"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"io"
	"net"
)

type Mail struct {
	Tls      bool        `json:"tls" yaml:"tls"`
	Server   string      `json:"server" yaml:"server"`
	From     FromAddress `json:"from" yaml:"from"`
	Username string      `json:"username" yaml:"username"`
	Password string      `json:"password" yaml:"password"`
}

func (m *Mail) loginInfo() sasl.Client {
	return sasl.NewPlainClient("", m.Username, m.Password)
}

func (m *Mail) mailCall(to []string, r io.Reader) error {
	host, _, err := net.SplitHostPort(m.Server)
	if err != nil {
		return err
	}
	if m.Tls {
		return smtp.SendMailTLS(m.Server, m.loginInfo(), m.From.String(), to, r)
	}
	if isLocalhost(host) {
		// internals of smtp.SendMail without STARTTLS for localhost testing
		dial, err := smtp.Dial(m.Server)
		if err != nil {
			return err
		}
		err = dial.Auth(m.loginInfo())
		if err != nil {
			return err
		}
		return dial.SendMail(m.From.String(), to, r)
	}
	return smtp.SendMail(m.Server, m.loginInfo(), m.From.String(), to, r)
}

func isLocalhost(host string) bool {
	// lookup host with resolver
	ip, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", host)
	if err != nil {
		return false
	}

	// missing list of addresses
	if len(ip) < 1 {
		return false
	}

	// if one ip is not loop back then this isn't localhost
	for _, i := range ip {
		if !i.IsLoopback() {
			return false
		}
	}

	// all addresses resolved are localhost
	return true
}

// SendMail sends a mail message to the provided mail addresses.
//
// The reader should follow the format of an RFC 822-style email.
//
// See github.com/emersion/go-smtp.SendMail
func (m *Mail) SendMail(to []*mail.Address, mailMessage io.Reader) error {
	// convert all to addresses to strings
	toStr := make([]string, len(to))
	for i := range toStr {
		toStr[i] = to[i].String()
	}

	return m.mailCall(toStr, mailMessage)
}
