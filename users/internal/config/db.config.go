package config

type DBConfig struct {
	Uri             string `mapstructure:"uri"`
	MaxConn         int    `mapstructure:"max_connection"`
	Minconn         int    `mapstructure:"min_connection"`
	MaxConnIdleTime int    `mapstructure:"max_connection_idle_time"`
}
