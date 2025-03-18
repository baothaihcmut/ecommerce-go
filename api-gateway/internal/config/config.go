package config



type CoreConfig struct{
	Server *ServerConfig `mapstructure:"server"`
	Address *AddressConfig `mapstructure:"addresses"`
	Web *WebConfig `mapstructure:"web"`
	Jwt *JwtConfig `mapstructure:"jwt"`
}