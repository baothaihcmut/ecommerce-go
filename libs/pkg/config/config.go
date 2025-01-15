package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(cfg interface{}, path string) (interface{}, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	// Parse the config into a struct
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Println("Error unmarshalling config:", err)
		return nil, err
	}
	return cfg, nil
}
