package simplemail

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"io"
)

type PreparedMail struct {
	simpleMail   *SimpleMail
	templateName string
	Header       mail.Header
	rcpt         []*mail.Address
	data         map[string]any
}

func (p *PreparedMail) SendMail() error {
	buf := new(bytes.Buffer)

	var bufHtml, bufTxt bytes.Buffer
	p.simpleMail.render(&bufHtml, &bufTxt, p.templateName, p.data)

	// setup html and text alternative headers
	var hHtml, hTxt mail.InlineHeader
	hHtml.Set("Content-Type", "text/html; charset=utf-8")
	hTxt.Set("Content-Type", "text/plain; charset=utf-8")

	createWriter, err := mail.CreateWriter(buf, p.Header)
	if err != nil {
		return err
	}
	inline, err := createWriter.CreateInline()
	if err != nil {
		return err
	}
	err = copyToPart(inline, hHtml, &bufHtml)
	if err != nil {
		return err
	}
	err = copyToPart(inline, hTxt, &bufTxt)
	if err != nil {
		return err
	}
	err = inline.Close()
	if err != nil {
		return err
	}

	return p.simpleMail.mailSender.SendMail(p.rcpt, buf)
}

func copyToPart(inline *mail.InlineWriter, header mail.InlineHeader, body io.Reader) error {
	part, err := inline.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, body)
	if err != nil {
		return err
	}
	return part.Close()
}
