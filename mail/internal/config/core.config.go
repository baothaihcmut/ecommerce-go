package config

type CoreConfig struct {
	Mailer   *MailerConfig   `mapstructure:"mailer"`
	RabbitMq *RabbitMqConfig `mapstructure:"rabbitmq"`
}
