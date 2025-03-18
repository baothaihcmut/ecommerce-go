package config


type WebConfig struct{
	Prefix string `mapstructure:"prefix"`
	Public []string `mapstructure:"public"`
}