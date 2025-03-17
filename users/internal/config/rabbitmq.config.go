package config

type RabbitMqConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Vhost    string `mapstructure:"vhost"`
	Username string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
}
