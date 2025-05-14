package main

import (
	"bytes"
	"text/template"

	"github.com/vanng822/go-premailer/premailer"
	gomail "gopkg.in/mail.v2"
)


type Mail struct {
	Domain 			string
	Host 			string
	Port 			int
	Username 		string
	Password 		string
	Encryption 		string
	FromAddress 	string
	FromName 		string
}


type Message struct {
	From 		string
	FromName 	string
	To 			[]string
	Subject 	string
	Attachments []string
	Data 		any
	DataMap 	map[string]any
}


func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	message := gomail.NewMessage()
	message.SetHeader("From", m.FromAddress)
	message.SetHeader("To", msg.To...)
	message.SetHeader("Subject", msg.Subject)
	message.SetBody("text/plain", plainMessage)
	message.AddAlternative("text/html", formattedMessage)


	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			message.Attach(x)
		}
	}

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}


func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}



func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}


func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses: false,
		CssToAttributes: false,
		KeepBangImportant: true,
	}
	
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}