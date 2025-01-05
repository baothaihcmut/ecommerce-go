package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig ServerConfig
	ConsulConfig ConsulConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ConsulConfig struct {
	Host string
	Port int
}

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Could not parse environment variable %s as int, using fallback %d", key, fallback)
	}
	return fallback
}

//	func getEnvBool(key string, fallback bool) bool {
//		if value, exists := os.LookupEnv(key); exists {
//			if boolValue, err := strconv.ParseBool(value); err == nil {
//				return boolValue
//			}
//			log.Printf("Warning: Could not parse environment variable %s as bool, using fallback %t", key, fallback)
//		}
//		return fallback
//	}
func LoadConfig(env string) (*Config, error) {
	if env == "dev" {
		godotenv.Load()
	}
	return &Config{
		ServerConfig: ServerConfig{
			Host: getEnvString("SERVER_HOST", "localhost"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		ConsulConfig: ConsulConfig{
			Host: getEnvString("CONSUL_HOST", "localhost"),
			Port: getEnvInt("CONSUL_PORT", 8500),
		},
	}, nil
}
