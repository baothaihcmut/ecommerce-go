package mailer

import (
	"bytes"
	"context"
	"html/template"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendMail(_ context.Context, arg SendMailArg) error
}

type MailerImpl struct {
	dialer     *gomail.Dialer
	mailConfig *config.MailerConfig
	logger 		logger.Logger 
}

type SendMailArg struct {
	Subject  string
	To       string
	Template string
	Data     any
}

func (g *MailerImpl) SendMail(ctx context.Context, arg SendMailArg) error {
	tmpl, err := template.ParseFiles("templates/" + arg.Template)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, arg.Data)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", g.mailConfig.Username)
	m.SetHeader("To", arg.To)
	m.SetHeader("Subject", arg.Subject)
	m.SetBody("text/html", body.String())
	if err := g.dialer.DialAndSend(m); err != nil {
		g.logger.Errorf(ctx,map[string]any{
			"to": arg.To,
			"subject": arg.Subject,
		}, "Error send mail: %v", err)
		return err
	}
	g.logger.Info(ctx,map[string]any{
		"to": arg.To,
		"subject": arg.Subject,
	}, "Send mail success")
	return nil
}
func NewMailer(dialer *gomail.Dialer, mailConfig *config.MailerConfig,logger logger.Logger) Mailer {
	return &MailerImpl{
		dialer:     dialer,
		mailConfig: mailConfig,
		logger: logger,
	}
}
