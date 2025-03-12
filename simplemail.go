package simplemail

import (
	"github.com/emersion/go-message/mail"
	htmlTemplate "html/template"
	"io"
	"io/fs"
	"log"
	textTemplate "text/template"
	"time"
)

type SimpleMail struct {
	mailSender    *Mail
	htmlTemplates *htmlTemplate.Template
	textTemplates *textTemplate.Template
}

func New(sender *Mail, templateFS fs.FS) (simpleMail *SimpleMail, err error) {
	m := &SimpleMail{mailSender: sender}
	m.htmlTemplates, err = htmlTemplate.New("mail").ParseFS(templateFS, "*.go.html")
	if err != nil {
		return
	}
	m.textTemplates, err = textTemplate.New("mail").ParseFS(templateFS, "*.go.txt")
	if err != nil {
		return
	}
	return m, nil
}

func (m *SimpleMail) render(wrHtml, wrTxt io.Writer, name string, data any) {
	err := m.htmlTemplates.ExecuteTemplate(wrHtml, name+".go.html", data)
	if err != nil {
		log.Printf("Failed to render mail html: %s: %s\n", name, err)
	}
	err = m.textTemplates.ExecuteTemplate(wrTxt, name+".go.txt", data)
	if err != nil {
		log.Printf("Failed to render mail text: %s: %s\n", name, err)
	}
}

// PrepareSingle constructs the headers for sending an email to the provided mail address.
func (m *SimpleMail) PrepareSingle(templateName, subject string, to *mail.Address, data map[string]any) *PreparedMail {
	return m.PrepareMany(templateName, subject, []*mail.Address{to}, data)
}

// PrepareMany constructs the headers for sending an email to the provided mail addresses.
func (m *SimpleMail) PrepareMany(templateName, subject string, to []*mail.Address, data map[string]any) *PreparedMail {
	p := &PreparedMail{
		simpleMail:   m,
		templateName: templateName,
		rcpt:         to,
		data:         data,
	}
	p.Header.SetDate(time.Now())
	p.Header.SetSubject(subject)
	p.Header.SetAddressList("From", []*mail.Address{m.mailSender.From.ToMailAddress()})
	p.Header.SetAddressList("To", to)
	p.Header.Set("Content-Type", "multipart/alternative")
	return p
}

// Send is a simplified version of PrepareSingle which sends the mail without allowing header modifications.
func (m *SimpleMail) Send(templateName, subject string, to *mail.Address, data map[string]any) error {
	p := m.PrepareSingle(templateName, subject, to, data)
	return p.SendMail()
}
