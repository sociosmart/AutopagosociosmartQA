package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"smartgas-payment/config"
	"text/template"
)

type SendMailOpts struct {
	TemplatePath string
	To           string
	Data         any
	Description  string
}

//go:generate mockery --name MailService --filename=mock_mail.go --inpackage=true
type MailService interface {
	SendMail(SendMailOpts) error
}

type mailService struct {
	config config.Config
}

func ProvideMailService(config config.Config) *mailService {
	return &mailService{
		config: config,
	}
}

func (ms *mailService) SendMail(opts SendMailOpts) error {

	tmpl, err := template.ParseFiles("templates/" + opts.TemplatePath)

	if err != nil {
		return err
	}

	var tmplBytes bytes.Buffer

	tmpl.Execute(&tmplBytes, opts.Data)

	auth := smtp.PlainAuth("", ms.config.SMTP.User, ms.config.SMTP.Password, ms.config.SMTP.Host)

	addr := fmt.Sprintf("%s:%v", ms.config.SMTP.Host, ms.config.SMTP.Port)

	to := []string{opts.To}

	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"+
		"\r\n"+
		"%s\r\n",
		opts.To,
		ms.config.FromEmail,
		opts.Description,
		tmplBytes.String(),
	)

	smtp.SendMail(addr, auth, ms.config.FromEmail, to, []byte(message))

	return nil

}
