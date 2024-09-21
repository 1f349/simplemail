package simplemail

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	htmlTemplate "html/template"
	"io"
	"io/fs"
	"log"
	textTemplate "text/template"
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

func (m *SimpleMail) Send(templateName, subject string, to *mail.Address, data map[string]any) error {
	var bufHtml, bufTxt bytes.Buffer
	m.render(&bufHtml, &bufTxt, templateName, data)
	return m.mailSender.SendMail(subject, []*mail.Address{to}, &bufHtml, &bufTxt)
}
