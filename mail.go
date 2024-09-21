package simplemail

import (
	"bytes"
	"context"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"io"
	"net"
	"time"
)

type Mail struct {
	Name     string      `json:"name" yaml:"name"`
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

func (m *Mail) SendMail(subject string, to []*mail.Address, htmlBody, textBody io.Reader) error {
	// generate the email in this template
	buf := new(bytes.Buffer)

	// setup mail headers
	var h mail.Header
	h.SetDate(time.Now())
	h.SetSubject(subject)
	h.SetAddressList("From", []*mail.Address{m.From.Address})
	h.SetAddressList("To", to)
	h.Set("Content-Type", "multipart/alternative")

	// setup html and text alternative headers
	var hHtml, hTxt mail.InlineHeader
	hHtml.Set("Content-Type", "text/html; charset=utf-8")
	hTxt.Set("Content-Type", "text/plain; charset=utf-8")

	createWriter, err := mail.CreateWriter(buf, h)
	if err != nil {
		return err
	}
	inline, err := createWriter.CreateInline()
	if err != nil {
		return err
	}
	partHtml, err := inline.CreatePart(hHtml)
	if err != nil {
		return err
	}
	if _, err := io.Copy(partHtml, htmlBody); err != nil {
		return err
	}
	partTxt, err := inline.CreatePart(hTxt)
	if err != nil {
		return err
	}
	if _, err := io.Copy(partTxt, textBody); err != nil {
		return err
	}

	// convert all to addresses to strings
	toStr := make([]string, len(to))
	for i := range toStr {
		toStr[i] = to[i].String()
	}

	return m.mailCall(toStr, buf)
}
