package config



type CoreConfig struct{
	Server *ServerConfig `mapstructure:"server"`
	Address *AddressConfig `mapstructure:"addresses"`
}