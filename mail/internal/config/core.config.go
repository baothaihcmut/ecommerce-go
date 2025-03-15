package config

type CoreConfig struct {
	Mailer *MailerConfig `mapstructure:"mailer"`
}
