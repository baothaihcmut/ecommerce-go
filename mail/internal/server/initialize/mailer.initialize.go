package initialize

import (
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"gopkg.in/gomail.v2"
)

func InitializeMailer(cfg *config.MailerConfig) (*gomail.Dialer, error) {
	return gomail.NewDialer(cfg.MailHost, cfg.MailPort, cfg.Username, cfg.Password), nil

}
