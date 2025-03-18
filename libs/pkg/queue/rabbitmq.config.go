package queue


type RabbitMqConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Vhost    string `mapstructure:"vhost"`
	Username string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	IsSecure bool `mapstructure:"is_secure"`
}
